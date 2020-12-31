import Vue from "vue";
import Vuex from "vuex";
import holidays from "@/api/holidays.js"
import { mapGetters } from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    lastPnL: 0.0,
    lastPl: 0.0,
    lastPd: 0.0,
    pdlc: 0.0,
    platformBuild: "",
    screen: {
      width: 0,
      height: 0,
      mode: "desktop"
    },
    sp: {},
    config: {
      isRiskVisible: false,
      refreshInterval: 3,
    },

    isHolidayToday: false,

    iStatusCounter: {
      0: 0, // ready
      1: 0, // running
      2: 0, // closed
      3: 0, // failed
      4: 0, // lost
      5: 0, // closing
    },
    holidays: {
      isHolidayToday: false,
      isUpdateRequired: false,
      userNotifyRequired: false,
      lastDay: "",
      dates: [],
      today: ""
    },

    iStatusColor: {
      0: "lightgreen",
      1: "white",
      2: "#9e9e9e",
      3: "tomato",
      4: "orange",
      5: "dodgerblue"
    },
  },
  mutations: {
    setHolidays(state, h) {
      state.holidays = h;
    },

    setHoliday(state, h) {
      state.isHolidayToday = h;
    },
    setPnL(state, n) {
      state.lastPnL = n;
    },
    setPd(state, n) {
      state.lastPd = n;
    },
    setPl(state, n) {
      state.lastPl = n;
    },
    setPdlc(state, n) {
      state.pdlc = n;
    },
    setSP(state, n) {
      state.sp = n;
    },

    setInstanceStatusCounter(state, c) {
      state.iStatusCounter = c;
    },
    setRefreshInterval(state, ri) {
      state.config.refreshInterval = ri;
      window.localStorage.setItem("refreshInterval", ri);
    },

    setRiskVisibility(state, b) {
      state.config.isRiskVisible = b;
      window.localStorage.setItem("riskVisible", b);
    },
    setScreen(state, s) {
      state.screen = s
    },
    setPlatformBuild(state, pb) {
      state.platformBuild = pb
    },
    setConfig(state, cfg) {
      window.localStorage.setItem("refreshInterval", cfg.refreshInterval);
      window.localStorage.setItem("riskVisible", cfg.isRiskVisible);
      state.config.isRiskVisible = cfg.isRiskVisible;
      state.config.refreshInterval = cfg.refreshInterval;
    }
  },
  getters: {
    holidays(state) {
      return state.holidays;
    },

    sp(state) {
      return state.sp;
    },


    holidayToday(state) {
      return state.isHolidayToday;
    },

    screen(state) {
      return state.screen;
    },
    platformBuild(state) {
      return state.platformBuild;
    },

    instanceStatusColor(state) {
      return state.iStatusColor;
    },

    instanceStatusCounter(state) {
      return state.iStatusCounter;
    },

    refreshInterval(state) {
      return state.config.refreshInterval;
    },

    isRiskVisible(state) {
      return state.config.isRiskVisible;
    },
    config(state) {
      return state.config;
    },

    isAtLeaseOneInstanceActive(state) {
      for (let k in state.iStatusCounter) {
        if (state.iStatusCounter[k] < 2) {
          return true;
        }
      }
      return false;
    },
  },
  actions: {
    // initHolidays(context) {
    //   let h = getHolidays();
    //   console.log(h)
    //   context.commit("setHolidays", h)
    // },
    init(context) {
      context.commit("setHolidays", holidays.getHolidays());
    }

  },
  modules: {
    // holidays: {
    //   state: {
    //     isHolidayToday: false,
    //     isUpdateRequired: false,
    //     lastDay: "",
    //     dates: [],
    //     today: ""
    //   },
    //   computed: {
    //     ...mapGetters([
    //       'isHolidayToday',
    //       'isUpdateRequired',
    //       'lastDay',
    //       'dates',
    //       'today'
    //     ])
    //   }
    // }
  },
});
