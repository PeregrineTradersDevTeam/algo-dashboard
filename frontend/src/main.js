import Vue from "vue";
import Vue2TouchEvents from "vue2-touch-events";
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";
import holidays from "@/api/holidays.js";

const defaultToastOptions = {
  position: "top-center",
  timeout: 2500,
  hideProgressBar: false,
  icon: true,
  draggable: true,
  pauseOnHover: true,
  transition: "Vue-Toastification__fade",
}

Vue.use(Toast, defaultToastOptions);

import VueConfirmDialog from 'vue-confirm-dialog'

Vue.use(VueConfirmDialog)

//import "buefy/dist/buefy.css";

import App from "./App.vue";

import router from "./router";
import store from "./store";
import { Dropdown, Dialog, ToastProgrammatic } from "buefy";

import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import screenProperties from "./mixins/screenProperties.js";

library.add(fas);

Vue.component("font-awesome-icon", FontAwesomeIcon);

Vue.config.productionTip = false;
//Vue.use(DropdownItem);
Vue.use(Dropdown);
Vue.use(Dialog);
Vue.use(ToastProgrammatic);

Vue.use({
  defaultIconComponent: "font-awesome-icon",
  defaultIconPack: "fa",
});

Vue.use(Vue2TouchEvents, {
  disableClick: false,
  touchClass: "",
  tapTolerance: 10,
  touchHoldTolerance: 400,
  swipeTolerance: 30,
  longTapTimeInterval: 400,
});



new Vue({
  router,
  store,
  render: (h) => h(App),
  created: function () {
    store.dispatch("init");
    //store.commit("setHolidays", holidays.getHolidays())

    let ri = 3;
    let sri = window.localStorage.getItem("refreshInterval");
    if (sri !== undefined && sri !== null) {
      ri = parseInt(sri);
      if (ri.IsNaN) {
        ri = 3;
      }
    }
    this.$store.commit("setRefreshInterval", ri);

    let rv = window.localStorage.getItem("riskVisible");
    if (rv) {
      this.$store.commit("setRiskVisibility", rv == "true");
    }

    let screen = screenProperties.get();
    this.$store.commit("setSP", screen)

    let that = this;
    window.addEventListener("resize", function (event) {
      that.$store.commit("setSP", screenProperties.get())
    });

  },

}).$mount("#app");
