<script setup>
import { computed, nextTick, onBeforeUnmount, ref } from 'vue'

const serverUrl = ref('ws://localhost:8080')
const connectionId = ref('room-1')
const sender = ref('Pedro')
const messageText = ref('Hola desde Vue')
const notificationRecipient = ref('user-123')
const notificationText = ref('Tenés una nueva notificación')

const socket = ref(null)
const status = ref('disconnected')
const events = ref([])
const messagesList = ref(null)

const isConnected = computed(() => status.value === 'connected')
const wsUrl = computed(() => `${serverUrl.value.replace(/\/$/, '')}/ws/${connectionId.value.trim()}`)

function addEvent(kind, payload) {
  events.value.push({
    id: crypto.randomUUID(),
    kind,
    payload,
    at: new Date().toLocaleTimeString(),
  })

  nextTick(() => {
    if (messagesList.value) {
      messagesList.value.scrollTop = messagesList.value.scrollHeight
    }
  })
}

function connect() {
  if (socket.value) {
    socket.value.close()
  }

  if (!connectionId.value.trim()) {
    addEvent('error', 'Tenés que definir una sala o usuario para conectarte.')
    return
  }

  status.value = 'connecting'
  addEvent('system', `Conectando a ${wsUrl.value}`)

  const ws = new WebSocket(wsUrl.value)
  socket.value = ws

  ws.onopen = () => {
    status.value = 'connected'
    addEvent('system', `Conectado como ${connectionId.value}`)
  }

  ws.onmessage = (event) => {
    try {
      addEvent('incoming', JSON.parse(event.data))
    } catch {
      addEvent('incoming', event.data)
    }
  }

  ws.onerror = () => {
    status.value = 'error'
    addEvent('error', 'Error de conexión WebSocket. Revisá que el backend esté corriendo.')
  }

  ws.onclose = () => {
    status.value = 'disconnected'
    addEvent('system', 'Conexión cerrada')
    socket.value = null
  }
}

function disconnect() {
  if (!socket.value) return
  socket.value.close()
}

function sendPayload(payload) {
  if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
    addEvent('error', 'Primero conectate al WebSocket.')
    return
  }

  socket.value.send(JSON.stringify(payload))
  addEvent('outgoing', payload)
}

function sendMessage() {
  const content = messageText.value.trim()
  if (!content) return

  sendPayload({
    type: 'message',
    id: connectionId.value.trim(),
    sender: sender.value.trim() || 'anonymous',
    content,
  })

  messageText.value = ''
}

function sendNotification() {
  const recipient = notificationRecipient.value.trim()
  const content = notificationText.value.trim()
  if (!recipient || !content) return

  sendPayload({
    type: 'notification',
    sender: sender.value.trim() || 'system',
    recipient,
    content,
  })
}

function clearEvents() {
  events.value = []
}

onBeforeUnmount(() => {
  if (socket.value) socket.value.close()
})
</script>

<template>
  <main class="page-shell">
    <section class="hero-card">
      <div>
        <p class="eyebrow">Vue + WebSocket</p>
        <h1>ChatBox Client</h1>
        <p class="subtitle">
          Cliente simple para probar el backend Go: salas, mensajes y notificaciones directas.
        </p>
      </div>

      <div class="status-pill" :class="status">
        <span class="dot"></span>
        {{ status }}
      </div>
    </section>

    <section class="grid">
      <aside class="panel controls-panel">
        <h2>Conexión</h2>

        <label>
          Backend WebSocket
          <input v-model="serverUrl" placeholder="ws://localhost:8080" />
        </label>

        <label>
          Sala o usuario conectado
          <input v-model="connectionId" placeholder="room-1" />
        </label>

        <label>
          Remitente
          <input v-model="sender" placeholder="Pedro" />
        </label>

        <div class="url-preview">
          <span>URL final</span>
          <code>{{ wsUrl }}</code>
        </div>

        <div class="button-row">
          <button class="primary" :disabled="status === 'connecting'" @click="connect">
            {{ isConnected ? 'Reconectar' : 'Conectar' }}
          </button>
          <button class="secondary" :disabled="!isConnected" @click="disconnect">
            Desconectar
          </button>
        </div>

        <div class="hint-box">
          Para probar chat, abrí dos pestañas conectadas a <strong>room-1</strong>.
          Para probar notificaciones, abrí otra pestaña conectada a <strong>user-123</strong>.
        </div>
      </aside>

      <section class="panel chat-panel">
        <div class="panel-header">
          <div>
            <h2>Eventos</h2>
            <p>Todo lo que entra y sale por WebSocket.</p>
          </div>
          <button class="ghost" @click="clearEvents">Limpiar</button>
        </div>

        <div ref="messagesList" class="messages-list">
          <div v-if="events.length === 0" class="empty-state">
            Todavía no hay eventos. Conectate y enviá un mensaje.
          </div>

          <article
            v-for="event in events"
            :key="event.id"
            class="message-card"
            :class="event.kind"
          >
            <div class="message-meta">
              <span>{{ event.kind }}</span>
              <time>{{ event.at }}</time>
            </div>
            <pre>{{ event.payload }}</pre>
          </article>
        </div>
      </section>
    </section>

    <section class="grid send-grid">
      <section class="panel form-panel">
        <h2>Enviar mensaje de sala</h2>
        <p>
          Usa <code>type: message</code> y <code>id</code> igual a la sala actual.
        </p>

        <textarea
          v-model="messageText"
          rows="4"
          placeholder="Escribí un mensaje para la sala actual"
          @keydown.ctrl.enter.prevent="sendMessage"
        ></textarea>

        <button class="primary" :disabled="!isConnected" @click="sendMessage">
          Enviar mensaje
        </button>
      </section>

      <section class="panel form-panel">
        <h2>Enviar notificación</h2>
        <p>
          Llega a quien esté conectado con ese ID como ruta WebSocket.
        </p>

        <label>
          Recipient
          <input v-model="notificationRecipient" placeholder="user-123" />
        </label>

        <textarea
          v-model="notificationText"
          rows="4"
          placeholder="Texto de la notificación"
          @keydown.ctrl.enter.prevent="sendNotification"
        ></textarea>

        <button class="primary" :disabled="!isConnected" @click="sendNotification">
          Enviar notificación
        </button>
      </section>
    </section>
  </main>
</template>
