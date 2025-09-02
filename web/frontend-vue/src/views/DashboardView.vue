<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import UptimeChart from '../components/UptimeChart.vue';

// State management
const sites = ref([]);
const newCheckTarget = ref('');
const errorMsg = ref(null);
const router = useRouter();
let socket = null;

// State untuk form
const newCheckType = ref('HTTP');
const newCheckKeyword = ref('');

// State untuk modal dan chart
const selectedSite = ref(null);
const chartData = ref(null);
const isLoadingChart = ref(false);
const activeRange = ref('24h');

const fetchSites = async () => {
  const token = localStorage.getItem('jwt_token');
  if (!token) {
    router.push('/login');
    return;
  }
  try {
    const response = await fetch('http://localhost:8080/api/sites', {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch sites');
    const data = await response.json();
    sites.value = (data || []).map(site => ({ ...site, last_checked: null, is_up: null }));
  } catch (err) {
    errorMsg.value = err.message;
  }
};

const handleAddSite = async () => {
  const token = localStorage.getItem('jwt_token');
  errorMsg.value = null;
  let payload = {
    check_type: newCheckType.value,
    check_target: newCheckTarget.value,
  };
  if (newCheckType.value === 'KEYWORD') {
    payload.check_keyword = newCheckKeyword.value;
  }
  try {
    const response = await fetch('http://localhost:8080/api/sites', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      const data = await response.json();
      throw new Error(data.error || 'Failed to add site');
    }
    newCheckTarget.value = '';
    newCheckKeyword.value = '';
    fetchSites();
  } catch (err) {
    errorMsg.value = err.message;
  }
};

const handleDeleteSite = async (siteId) => {
  const token = localStorage.getItem('jwt_token');
  if (!confirm('Are you sure you want to delete this monitor?')) return;
  try {
    const response = await fetch(`http://localhost:8080/api/sites/${siteId}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to delete site');
    sites.value = sites.value.filter(s => s.id !== siteId);
  } catch (err) {
    errorMsg.value = err.message;
  }
};

const connectWebSocket = () => {
  const token = localStorage.getItem('jwt_token');
  socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
  socket.onopen = () => console.log('WebSocket connection established.');
  socket.onmessage = (event) => {
    const update = JSON.parse(event.data);
    const siteIndex = sites.value.findIndex(s => s.id === update.site_id);
    if (siteIndex !== -1) {
      sites.value[siteIndex].is_up = update.is_up;
      sites.value[siteIndex].response_time_ms = update.response_time_ms;
      sites.value[siteIndex].status_code = update.status_code;
      sites.value[siteIndex].last_checked = new Date(update.checked_at);
    }
  };
  socket.onclose = () => console.log('WebSocket connection closed.');
  socket.onerror = (error) => console.error('WebSocket error:', error);
};

const handleLogout = () => {
  localStorage.removeItem('jwt_token');
  router.push('/login');
};

const showHistory = (site) => {
  selectedSite.value = site;
  activeRange.value = '24h';
  fetchHistoryData(site.id, '24h');
};

const fetchHistoryData = async (siteId, range) => {
  isLoadingChart.value = true;
  chartData.value = null;
  try {
    const token = localStorage.getItem('jwt_token');
    const response = await fetch(`http://localhost:8080/api/sites/${siteId}/history?range=${range}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    if (!response.ok) throw new Error('Failed to fetch history');
    const historyData = await response.json();
    chartData.value = {
      labels: historyData.map(d => new Date(d.checked_at)),
      datasets: [{
        label: 'Response Time (ms)',
        backgroundColor: '#3b82f6',
        borderColor: '#3b82f6',
        data: historyData.map(d => d.response_time_ms),
        tension: 0.1,
        pointRadius: 2,
      }]
    };
  } catch (err) {
    console.error(err);
  } finally {
    isLoadingChart.value = false;
  }
};

const changeRange = (newRange) => {
  if (!selectedSite.value) return;
  activeRange.value = newRange;
  fetchHistoryData(selectedSite.value.id, newRange);
};

const closeHistoryModal = () => {
  selectedSite.value = null;
  chartData.value = null;
};

onMounted(() => {
  fetchSites();
  connectWebSocket();
});

onUnmounted(() => {
  if (socket) socket.close();
});

const formatDate = (date) => {
  if (!date) return 'Checking...';
  return new Date(date).toLocaleString();
};

const getStatusColor = (status) => {
  if (status === true) return 'bg-green-500';
  if (status === false) return 'bg-red-500';
  return 'bg-gray-400';
};

const getStatusText = (status) => {
  if (status === true) return 'Online';
  if (status === false) return 'Offline';
  return 'Pending';
};
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header yang diperbaiki -->
    <nav class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
          <div class="flex items-center">
            <div class="w-8 h-8 bg-blue-600 rounded-md flex items-center justify-center mr-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <h1 class="text-xl font-semibold text-gray-900">Go-Pulse</h1>
          </div>
          <div class="flex items-center">
            <button 
              @click="handleLogout" 
              class="flex items-center text-sm text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md hover:bg-gray-100 transition-colors"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Konten utama -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Form untuk menambah monitor -->
      <div class="px-4 py-6 sm:px-0">
        <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h2 class="text-lg font-medium text-gray-900 mb-4">Add a new monitor</h2>
          <form @submit.prevent="handleAddSite" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div class="md:col-span-2">
                <label for="check_target" class="block text-sm font-medium text-gray-700 mb-1">Target URL</label>
                <input 
                  v-model="newCheckTarget" 
                  type="text" 
                  :placeholder="newCheckType === 'TCP' ? 'example.com:443' : 'https://example.com'" 
                  required 
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                >
              </div>
              <div>
                <label for="check_type" class="block text-sm font-medium text-gray-700 mb-1">Check Type</label>
                <select 
                  v-model="newCheckType" 
                  id="check_type" 
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                >
                  <option>HTTP</option>
                  <option>KEYWORD</option>
                  <option>TCP</option>
                </select>
              </div>
            </div>
            
            <div v-if="newCheckType === 'KEYWORD'">
              <label for="check_keyword" class="block text-sm font-medium text-gray-700 mb-1">Keyword to find</label>
              <input 
                v-model="newCheckKeyword" 
                type="text" 
                placeholder="e.g., Copyright" 
                required 
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              >
            </div>
            
            <p v-if="errorMsg" class="text-red-500 text-sm">{{ errorMsg }}</p>
            
            <button 
              type="submit" 
              class="px-4 py-2 font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors"
            >
              Add Monitor
            </button>
          </form>
        </div>
      </div>

      <!-- Daftar monitor -->
      <div class="px-4 py-6 sm:px-0">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Your Monitors</h2>
        
        <div v-if="sites.length === 0" class="text-center text-gray-500 py-8 bg-white rounded-lg border border-gray-200">
          You are not monitoring any sites yet.
        </div>
        
        <div v-else class="grid grid-cols-1 gap-4">
          <div 
            v-for="site in sites" 
            :key="site.id" 
            class="bg-white p-5 rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <span class="flex h-3 w-3 relative">
                  <span v-if="site.is_up === true" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span v-else-if="site.is_up === false" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                  <span :class="{
                    'bg-green-500': site.is_up === true, 
                    'bg-red-500': site.is_up === false, 
                    'bg-gray-400': site.is_up === null
                  }" class="relative inline-flex rounded-full h-3 w-3"></span>
                </span>
                <div>
                  <h3 class="font-medium text-gray-900">{{ site.check_target }}</h3>
                  <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full">{{ site.check_type }}</span>
                </div>
              </div>
              
              <div class="flex items-center space-x-4">
                <div class="text-sm text-gray-600 text-right">
                  <div v-if="site.status_code">Status: {{ site.status_code }}</div>
                  <div v-if="site.response_time_ms">Response: {{ site.response_time_ms }} ms</div>
                  <div>Last Checked: {{ formatDate(site.last_checked) }}</div>
                </div>
                <button 
                  @click="handleDeleteSite(site.id)" 
                  class="text-red-500 hover:text-red-700 transition-colors"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>