import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { fetchChats } from "./chatSlice";
import type { RootState, AppDispatch } from "../../store";

export default function ChatsPage() {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { chats, loading, error } = useSelector(
    (state: RootState) => state.chat
  );

  useEffect(() => {
    dispatch(fetchChats());
  }, [dispatch]);

  const handleLogout = () => {
    // dispatch(logout());
    // navigate('/login');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-gray-600">Loading chats...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center py-4">
          <h1 className="text-2xl font-bold text-gray-900">Chats</h1>
          <button
            onClick={handleLogout}
            className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
          >
            Logout
          </button>
        </div>

        {error && (
          <div className="mb-4 bg-red-50 border border-red-400 text-red-700 px-4 py-3 rounded relative">
            {error}
          </div>
        )}

        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {chats.map((chat) => (
              <li key={chat.id}>
                <div className="px-4 py-4 flex items-center sm:px-6">
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center justify-between">
                      <p className="text-sm font-medium text-indigo-600 truncate">
                        {chat.title}
                      </p>
                      <div className="ml-2 flex-shrink-0 flex">
                        <p className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                          {new Date(chat.updatedAt).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                    {chat.lastMessage && (
                      <div className="mt-2">
                        <p className="text-sm text-gray-500">
                          {chat.lastMessage}
                        </p>
                      </div>
                    )}
                  </div>
                </div>
              </li>
            ))}
            {chats.length === 0 && !loading && (
              <li className="px-4 py-4 sm:px-6">
                <div className="text-center text-gray-500">
                  No chats available
                </div>
              </li>
            )}
          </ul>
        </div>
      </div>
    </div>
  );
}
