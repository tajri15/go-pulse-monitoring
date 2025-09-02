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
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex-shrink-0 flex items-center">
            <h1 class="text-xl font-bold text-blue-600">Go-Pulse</h1>
          </div>
          <div class="flex items-center">
            <button @click="handleLogout" class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-100">
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <div class="bg-white p-6 rounded-lg shadow">
          <h2 class="text-lg font-medium mb-4">Add a new monitor</h2>
          <form @submit.prevent="handleAddSite" class="space-y-4">
            <div class="flex flex-col md:flex-row space-y-4 md:space-y-0 md:space-x-4">
              <div class="flex-grow">
                <label for="check_target" class="block text-sm font-medium">Target</label>
                <input v-model="newCheckTarget" type="text" :placeholder="newCheckType === 'TCP' ? 'example.com:443' : 'https://example.com'" required class="w-full px-3 py-2 mt-1 border rounded-md">
              </div>
              <div>
                <label for="check_type" class="block text-sm font-medium">Type</label>
                <select v-model="newCheckType" id="check_type" class="w-full px-3 py-2 mt-1 border rounded-md h-full">
                  <option>HTTP</option>
                  <option>KEYWORD</option>
                  <option>TCP</option>
                </select>
              </div>
            </div>
            <div v-if="newCheckType === 'KEYWORD'">
              <label for="check_keyword" class="block text-sm font-medium">Keyword to find</label>
              <input v-model="newCheckKeyword" type="text" placeholder="e.g., Copyright" required class="w-full px-3 py-2 mt-1 border rounded-md">
            </div>
            <p v-if="errorMsg" class="text-red-500 text-sm">{{ errorMsg }}</p>
            <button type="submit" class="px-4 py-2 font-bold text-white bg-blue-600 rounded-md hover:bg-blue-700">Add Monitor</button>
          </form>
        </div>
      </div>

      <div class="px-4 py-6 sm:px-0">
        <div class="space-y-4">
          <div v-if="sites.length === 0" class="text-center text-gray-500 py-8">
            You are not monitoring any sites yet.
          </div>
          <div v-else v-for="site in sites" :key="site.id" class="bg-white p-4 rounded-lg shadow flex items-center justify-between flex-wrap">
            <div class="flex items-center space-x-4 mb-2 md:mb-0">
              <span class="flex h-3 w-3 relative">
                <span v-if="site.is_up === true" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                <span v-else-if="site.is_up === false" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                <span :class="{'bg-green-500': site.is_up === true, 'bg-red-500': site.is_up === false, 'bg-gray-400': site.is_up === null}" class="relative inline-flex rounded-full h-3 w-3"></span>
              </span>
              <button @click="showHistory(site)" class="font-medium text-blue-600 hover:underline text-left">{{ site.check_target }}</button>
              <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full">{{ site.check_type }}</span>
            </div>
            <div class="flex items-center space-x-6 text-sm text-gray-600">
              <span>{{ site.is_up === null ? 'Pending' : `Status: ${site.status_code}` }}</span>
              <span>{{ site.response_time_ms === undefined ? '' : `Response: ${site.response_time_ms} ms` }}</span>
              <span class="min-w-[190px]">Last Checked: {{ formatDate(site.last_checked) }}</span>
              <button @click="handleDeleteSite(site.id)" class="text-red-500 hover:text-red-700 font-semibold">Delete</button>
            </div>
          </div>
        </div>
      </div>
    </main>

    <div v-if="selectedSite" @click.self="closeHistoryModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex justify-center items-center p-4">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-4xl h-full max-h-[70vh] flex flex-col">
        <div class="p-4 border-b flex justify-between items-center">
          <h3 class="text-lg font-medium">History: {{ selectedSite.check_target }}</h3>
          <div class="flex space-x-2">
            <button @click="changeRange('24h')" :class="{'bg-blue-500 text-white': activeRange === '24h', 'bg-gray-200': activeRange !== '24h'}" class="px-3 py-1 text-sm rounded-md">24 Hours</button>
            <button @click="changeRange('7d')" :class="{'bg-blue-500 text-white': activeRange === '7d', 'bg-gray-200': activeRange !== '7d'}" class="px-3 py-1 text-sm rounded-md">7 Days</button>
            <button @click="changeRange('30d')" :class="{'bg-blue-500 text-white': activeRange === '30d', 'bg-gray-200': activeRange !== '30d'}" class="px-3 py-1 text-sm rounded-md">30 Days</button>
          </div>
          <button @click="closeHistoryModal" class="text-2xl leading-none text-gray-500 hover:text-gray-800">&times;</button>
        </div>
        <div class="p-4 flex-grow relative">
          <div v-if="isLoadingChart" class="absolute inset-0 flex items-center justify-center text-gray-500">Loading chart data...</div>
          <UptimeChart v-if="chartData" :chart-data="chartData" />
        </div>
      </div>
    </div>
  </div>
</template>