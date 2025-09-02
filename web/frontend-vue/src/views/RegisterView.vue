<script setup>
import { ref } from 'vue';
import { useRouter, RouterLink } from 'vue-router';

const username = ref('');
const email = ref('');
const password = ref('');
const errorMsg = ref(null);
const successMsg = ref(null);
const router = useRouter();

const handleRegister = async () => {
  errorMsg.value = null;
  successMsg.value = null;
  try {
    const response = await fetch('http://localhost:8080/api/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        email: email.value,
        password: password.value,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Registration failed');
    }

    // Tampilkan pesan sukses dan arahkan ke halaman login setelah beberapa saat
    successMsg.value = 'Registration successful! Redirecting to login...';
    setTimeout(() => {
      router.push('/login');
    }, 2000);

  } catch (err) {
    errorMsg.value = err.message;
  }
};
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
      <h1 class="text-2xl font-bold text-center">Create Go-Pulse Account</h1>
      <form @submit.prevent="handleRegister" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium">Username</label>
          <input v-model="username" type="text" id="username" required class="w-full px-3 py-2 mt-1 border rounded-md focus:outline-none focus:ring focus:ring-blue-200">
        </div>
        <div>
          <label for="email" class="block text-sm font-medium">Email</label>
          <input v-model="email" type="email" id="email" required class="w-full px-3 py-2 mt-1 border rounded-md focus:outline-none focus:ring focus:ring-blue-200">
        </div>
        <div>
          <label for="password" class="block text-sm font-medium">Password</label>
          <input v-model="password" type="password" id="password" required minlength="6" class="w-full px-3 py-2 mt-1 border rounded-md focus:outline-none focus:ring focus:ring-blue-200">
        </div>
        <div v-if="errorMsg" class="text-red-500 text-sm text-center">
          {{ errorMsg }}
        </div>
        <div v-if="successMsg" class="text-green-500 text-sm text-center">
          {{ successMsg }}
        </div>
        <button type="submit" class="w-full px-4 py-2 font-bold text-white bg-blue-600 rounded-md hover:bg-blue-700">
          Register
        </button>
      </form>
       <div class="text-center mt-4">
        <p class="text-sm">
          Sudah punya akun? 
          <RouterLink to="/login" class="font-medium text-blue-600 hover:underline">Login di sini</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>