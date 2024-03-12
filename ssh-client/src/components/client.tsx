import { defineComponent, ref, onMounted } from "vue";
import { Terminal } from 'xterm';
import { AttachAddon } from 'xterm-addon-attach';
import { FitAddon } from 'xterm-addon-fit';
import "xterm/css/xterm.css"

export default defineComponent({
  setup() {
    const tm = ref<HTMLElement>()
    onMounted(() => {
      const terminal = new Terminal({
        rows: 45,
      });
      const fitAddon = new FitAddon();
      terminal.loadAddon(fitAddon);
      // const tm = document.getElementById('terminal')
      if (tm.value) {
        terminal.open(tm.value);
        fitAddon.fit();

        let webSocket = new WebSocket('ws://localhost:8080/web-socket/ssh');

        const sendSize = () => {
          const windowSize = { high: terminal.rows, width: terminal.cols };
          const blob = new Blob([JSON.stringify(windowSize)], { type: 'application/json' });
          webSocket.send(blob);
        }

        webSocket.onopen = sendSize;

        const resizeScreen = () => {
          fitAddon.fit();
          sendSize();
        }
        window.addEventListener('resize', resizeScreen, false);

        const attachAddon = new AttachAddon(webSocket);
        terminal.loadAddon(attachAddon);
      }
    })

    return () => <>
      <div ref={tm} ></div>
    </>
  }
})
