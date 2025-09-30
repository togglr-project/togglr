import { WSClient } from './ws';
import { handleEvent, type RealtimeEvent } from './handlers';

export type InitOptions = {
  projectId: string;
  envId: string; // numeric string
  token?: string;
  onStatus?: (status: 'connected' | 'disconnected' | 'error') => void;
};

let client: WSClient | null = null;

export function initRealtime(opts: InitOptions) {
  console.log('[Realtime] Initializing with options:', opts);
  
  // Close previous if any
  client?.stop();

  client = new WSClient({
    projectId: opts.projectId,
    envId: opts.envId,
    token: opts.token,
    onOpen: () => {
      console.log('[Realtime] Connection opened');
      opts.onStatus?.('connected');
    },
    onClose: (ev) => {
      console.log('[Realtime] Connection closed', ev);
      opts.onStatus?.('disconnected');
    },
    onError: (ev) => {
      console.error('[Realtime] Connection error', ev);
      opts.onStatus?.('error');
    },
    onMessage: (evt: RealtimeEvent) => {
      console.log('[Realtime] Event received:', evt);
      const qc = typeof window !== 'undefined' ? window.__RQ : undefined;
      handleEvent(qc, evt);
    },
  });

  client.start();

  return () => client?.stop();
}

export function stopRealtime() {
  client?.stop();
  client = null;
}
