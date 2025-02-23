import logging
import sys
from pathlib import Path
import time

import psycopg2
from alembic import command
from alembic.config import Config
from liive_sql_models.config import get_database_url

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def wait_for_db(max_retries=10, delay=2):
    """Wait for database to be ready."""
    url = get_database_url()
    retries = 0

    while retries < max_retries:
        try:
            conn_params = {
                "dbname": url.database,
                "host": url.host,
                "port": url.port or 5432,
                "user": url.username,
                "password": url.password,
            }

            logger.info(
                "Attempting to connect to PostgreSQL with params: host=%s port=%s user=%s dbname=%s",
                conn_params["host"],
                conn_params["port"],
                conn_params["user"],
                conn_params["dbname"],
            )

            with psycopg2.connect(**conn_params) as conn:
                with conn.cursor() as cur:
                    cur.execute("SELECT version()")
                    version = cur.fetchone()[0]
                    logger.info("Successfully connected to database: %s", version)
                    return True
        except psycopg2.Error as e:
            logger.warning(
                "Failed to connect to database (attempt %d/%d): %s", retries + 1, max_retries, e
            )
            if "database" in str(e) and "does not exist" in str(e):
                # Try connecting to default database to create our database
                try:
                    temp_params = conn_params.copy()
                    temp_params["dbname"] = "postgres"
                    with psycopg2.connect(**temp_params) as conn:
                        conn.autocommit = True
                        with conn.cursor() as cur:
                            cur.execute(f"CREATE DATABASE {conn_params['dbname']}")
                            logger.info("Created database %s", conn_params["dbname"])
                            return True
                except psycopg2.Error as create_err:
                    logger.warning("Failed to create database: %s", create_err)
            retries += 1
            if retries < max_retries:
                time.sleep(delay)

    logger.error("Could not connect to database after %d attempts", max_retries)
    return False


def run_migrations() -> None:
    """Run all pending migrations."""
    try:
        # Get the directory containing migrations
        migrations_dir = Path(__file__).parent / "migrations"

        # Create Alembic configuration
        alembic_cfg = Config()
        alembic_cfg.set_main_option("script_location", str(migrations_dir))
        alembic_cfg.set_main_option("sqlalchemy.url", str(get_database_url()))

        # Run migrations
        command.upgrade(alembic_cfg, "head")
        logger.info("Successfully ran all migrations")
    except Exception as e:
        logger.error("Failed to run migrations: %s", e)
        sys.exit(1)


def main() -> None:
    """Main entry point for the migrator."""
    logger.info("Starting database migration process")

    if not wait_for_db():
        logger.error("Failed to connect to database")
        sys.exit(1)

    run_migrations()
    logger.info("Database migration process completed")


if __name__ == "__main__":
    main()
