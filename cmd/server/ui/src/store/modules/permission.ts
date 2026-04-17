import { defineStore } from "pinia";
import { ref } from "vue";
import type { RouteRecordRaw } from "vue-router";

/**
 * 预留：业务系统里常由后端 /api/menus/routes 下发菜单，
 * 再用 Vite 的 import.meta.glob 解析 views 下组件路径并 addRoute。
 * 本控制台当前为静态路由，仅保留结构与占位状态。
 */
export const usePermissionStore = defineStore("permission", () => {
  const dynamicRoutesReady = ref(true);
  const menuRoutes = ref<RouteRecordRaw[]>([]);

  function setMenuRoutes(routes: RouteRecordRaw[]) {
    menuRoutes.value = routes;
  }

  return { dynamicRoutesReady, menuRoutes, setMenuRoutes };
});
