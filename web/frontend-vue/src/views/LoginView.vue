<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const email = ref('');
const password = ref('');
const errorMsg = ref(null);
const router = useRouter();

const handleLogin = async () => {
  errorMsg.value = null;
  try {
    const response = await fetch('http://localhost:8080/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email.value,
        password: password.value,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Login failed');
    }

    // Simpan token ke localStorage
    localStorage.setItem('jwt_token', data.token);
    // Arahkan ke dashboard
    router.push('/');

  } catch (err) {
    errorMsg.value = err.message;
  }
};
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
      <h1 class="text-2xl font-bold text-center">Login to Go-Pulse</h1>
      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label for="email" class="block text-sm font-medium">Email</label>
          <input v-model="email" type="email" id="email" required class="w-full px-3 py-2 mt-1 border rounded-md focus:outline-none focus:ring focus:ring-blue-200">
        </div>
        <div>
          <label for="password" class="block text-sm font-medium">Password</label>
          <input v-model="password" type="password" id="password" required class="w-full px-3 py-2 mt-1 border rounded-md focus:outline-none focus:ring focus:ring-blue-200">
        </div>
        <div v-if="errorMsg" class="text-red-500 text-sm text-center">
          {{ errorMsg }}
        </div>
        <button type="submit" class="w-full px-4 py-2 font-bold text-white bg-blue-600 rounded-md hover:bg-blue-700">
          Login
        </button>
      </form>
    </div>
  </div>
</template>