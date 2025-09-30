import { apiConfiguration } from '../api/apiClient';

export type WSOptions = {
  projectId: string;
  envId: string; // numeric string
  token?: string; // optional Bearer token
  onOpen?: () => void;
  onClose?: (ev: CloseEvent) => void;
  onError?: (ev: Event) => void;
  onMessage?: (data: any) => void;
};

// Build WS base URL from runtime config
const getWSBaseURL = (): string => {
  if (typeof window !== 'undefined' && window.TOGGLR_CONFIG?.WS_BASE_URL) {
    const raw = window.TOGGLR_CONFIG.WS_BASE_URL;
    
    // If already has ws:// or wss://, use as is
    if (raw.startsWith('ws://') || raw.startsWith('wss://')) {
      return raw.replace(/\/$/, '');
    }
    
    try {
      const u = new URL(raw);
      if (u.protocol === 'http:') u.protocol = 'ws:';
      if (u.protocol === 'https:') u.protocol = 'wss:';
      return u.toString().replace(/\/$/, '');
    } catch {
      // If it's a bare host like localhost:8082, prefix with ws(s) based on location
      const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      return `${proto}//${raw}`;
    }
  }
  
  console.log('[Realtime] WS_BASE_URL not found, inferring from API base');
  // If not specified, try to infer from API base: replace protocol/port
  try {
    const api = new URL(apiConfiguration.basePath || window.location.origin);
    const proto = api.protocol === 'https:' ? 'wss:' : 'ws:';
    return `${proto}//${api.host}`;
  } catch {
    return (window.location.protocol === 'https:' ? 'wss:' : 'ws:') + '//' + window.location.host;
  }
};

export class WSClient {
  private opts: WSOptions;
  private ws: WebSocket | null = null;
  private backoff = 1000; // start 1s
  private stopped = false;
  private pingInterval: number | null = null;

  constructor(opts: WSOptions) {
    this.opts = opts;
  }

  start() {
    this.stopped = false;
    this.connect();
  }

  stop() {
    this.stopped = true;
    this.stopPingPong();
    if (this.ws) {
      try { this.ws.close(); } catch {}
      this.ws = null;
    }
  }

  private connect() {
    const base = getWSBaseURL();
    
    // Convert HTTP URL to WebSocket URL
    const wsBase = base.replace('http://', 'ws://').replace('https://', 'wss://');
    console.log('[Realtime] WS base URL converted:', wsBase);
    
    const url = new URL('/api/ws', wsBase);
    url.searchParams.set('project_id', this.opts.projectId);
    url.searchParams.set('env_id', this.opts.envId);
    
    // Pass token through query parameter instead of subprotocol
    if (this.opts.token) {
      url.searchParams.set('token', this.opts.token);
      console.log('[Realtime] Using token in query parameter');
    } else {
      console.log('[Realtime] No token provided');
    }

    // No subprotocols needed
    const protocols: string[] = [];

    const finalURL = url.toString();
    try {
      this.ws = new WebSocket(finalURL, protocols);
    } catch (e) {
      console.error('[Realtime] WS connect error', e);
      this.scheduleReconnect();
      return;
    }

    this.ws.onopen = () => {
      console.log('[Realtime] WS connected successfully');
      this.backoff = 1000;
      this.opts.onOpen?.();
      
      // Start ping/pong to keep connection alive
      this.startPingPong();
    };

    this.ws.onclose = (ev) => {
      console.log('[Realtime] WS closed', { code: ev.code, reason: ev.reason, wasClean: ev.wasClean });
      this.opts.onClose?.(ev);
      if (!this.stopped) {
        console.log('[Realtime] WS scheduling reconnect...');
        this.scheduleReconnect();
      }
    };

    this.ws.onerror = (ev) => {
      console.error('[Realtime] WS error', ev);
      this.opts.onError?.(ev);
    };

    this.ws.onmessage = (ev) => {
      try {
        const data = JSON.parse(ev.data as string);
        console.log('[Realtime] WS message received', data);
        this.opts.onMessage?.(data);
      } catch (e) {
        console.warn('[Realtime] WS message parse error', e);
      }
    };
  }

  private scheduleReconnect() {
    const delay = Math.min(this.backoff, 30000); // Max 30 seconds
    console.log(`[Realtime] WS reconnecting in ${delay}ms...`);
    setTimeout(() => {
      if (!this.stopped) {
        this.backoff = Math.min(this.backoff * 1.5, 30000); // Slower backoff
        console.log('[Realtime] WS attempting reconnect...');
        this.connect();
      }
    }, delay);
  }

  private startPingPong() {
    this.stopPingPong(); // Clear any existing interval
    
    // Send ping every 30 seconds using WebSocket ping frame
    this.pingInterval = window.setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        try {
          // Use WebSocket ping frame if available, otherwise send ping message
          if ('ping' in this.ws && typeof this.ws.ping === 'function') {
            this.ws.ping();
          } else {
            // Fallback: send ping message
            this.ws.send(JSON.stringify({ type: 'ping' }));
          }
        } catch (e) {
          console.warn('[Realtime] Ping failed:', e);
        }
      }
    }, 30000);
  }

  private stopPingPong() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }
}
