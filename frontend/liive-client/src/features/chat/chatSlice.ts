import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import api from '../../services/api';

interface Chat {
  id: string;
  title: string;
  lastMessage?: string;
  updatedAt: string;
}

interface ChatState {
  chats: Chat[];
  loading: boolean;
  error: string | null;
}

const initialState: ChatState = {
  chats: [],
  loading: false,
  error: null,
};

export const fetchChats = createAsyncThunk('chat/fetchChats', async () => {
  const response = await api.chat.get('/chats');
  return response.data;
});

const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchChats.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchChats.fulfilled, (state, action) => {
        state.loading = false;
        state.chats = action.payload;
      })
      .addCase(fetchChats.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch chats';
      });
  },
});

export const { clearError } = chatSlice.actions;
export default chatSlice.reducer; 