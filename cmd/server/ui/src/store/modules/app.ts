import { defineStore } from "pinia";
import { ref } from "vue";

export const useAppStore = defineStore("app", () => {
  const appTitle = ref("Event Bus 控制台");
  return { appTitle };
});
