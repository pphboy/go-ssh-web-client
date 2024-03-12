import { defineComponent } from "vue";
import Client from "./components/client";


export default defineComponent({
  setup() {
    return () => <>
      <Client></Client>
    </>
  }
})
