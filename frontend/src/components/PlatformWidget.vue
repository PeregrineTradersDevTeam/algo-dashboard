<template>
  <div
    class="cell"
    v-bind:style="{ 'background-color': backcolor }"
    @mouseover="hover = true"
    @mouseleave="hover = false"
    :class="{ cellhovered: hover }"
  >
    <div>
      <span>{{ title }}</span>
      <Dropdown
        :items="menu"
        @onClick="restartPlatform"
        style="display:inline;"
      />
    </div>
    <!-- <span>{{ title }}</span>
      <b-dropdown aria-role="list">
        <font-awesome-icon icon="caret-down" slot="trigger" />
        <b-dropdown-item
          :disabled="restartDisabled"
          aria-role="listitem"
          v-on:click="restartPlatform"
          >Restart Platform</b-dropdown-item
        >
      </b-dropdown> -->
    <div class="connector" v-for="c in platform.connectors" v-bind:key="c.name">
      <span v-bind:style="{ color: config.connectorStatusColor[c.status] }"
        >[{{ config.connectorStatusSymbol[c.status] }}] {{ c.name }}</span
      >
    </div>
    <div class="widget-footer">
      {{ platform.server_time }}
    </div>
    <div class="log-button" v-on:click="showModal('PTF')">
      {{ platform.log_items_count }}
    </div>
    <PlatformModal ref="modal" />
  </div>
</template>

<script>
import PlatformModal from "@/components/PlatformModal.vue";
import { ToastProgrammatic as Toast } from "buefy";
import Dropdown from "@/components/Dropdown.vue";

export default {
  name: "PlatformWidget",
  components: {
    PlatformModal: PlatformModal,
    Dropdown,
  },
  props: {
    title: String,
    isStopped: Boolean,
  },
  data: function() {
    return {
      restartDisabled: false,
      menu: ["Restart Platform"],
      selected: null,
      hover: false,
      platform: {
        connectors: [],
        log_items_count: 0,
        status: 0,
        status_at: 0,
        pnl: 0,
        pd: 0,
        pl: 0,
        server_time: "",
      },
      logButtonColor: "white",
      config: {
        // connector status name abbrevation
        connectorStatusSymbol: {
          0: "I", // INSTANTIATED
          1: "S", // STARTED
          2: "R", // READY
          3: "D", // DOWN
          4: "K", // KILLED
          5: "E", // ERROR
        },
        // connector label color
        connectorStatusColor: {
          0: "green",
          1: "green",
          2: "green",
          3: "red",
          4: "red",
          5: "red",
        },

        // widget background color
        platformStatusColor: {
          0: "lightgreen", // READY
          1: "#23d160", // RUNNING
          2: "silver", // CLOSED
          3: "tomato", // FALIED
        },
      },
    };
  },

  computed: {
    backcolor: function() {
      if (this.isStopped) {
        return this.config.platformStatusColor[2];
      }
      let failed = 0;
      //let res = "dimgray"; // back color based on global platform status
      //let cres = "dimgray"  // back color based on connectors statuses

      let res = this.config.platformStatusColor[this.platform.status];
      if (this.platform.status > 1) {
        return res;
      }

      if (
        this.platform.connectors == null ||
        this.platform.connectors.length == 0
      ) {
        return res;
      }

      for (let i = 0; i < this.platform.connectors.length; i++) {
        if (this.platform.connectors[i].status > 2) {
          failed++;
        }
      }
      switch (failed) {
        case 0:
          break;
        case 1:
          res = "orange";
          break;
        default:
          res = this.config.platformStatusColor[3];
          break;
      }
      return res;
    },
  },
  methods: {
    restartPlatform: function() {
      this.$buefy.dialog.confirm({
        message: "Restart Platform?",
        focusOn: "cancel",
        animation: "",
        type: "is-danger",
        onConfirm: () => this.sendRestartPlatformRequest(),
      });
    },
    sendRestartPlatformRequest: function() {
      fetch("/api/platform/restart", {
        method: "POST",
      })
        .then((r) => r.json())
        .then(function(j) {
          if (j.received === true) {
            Toast.open({
              message: "restart request received by server!",
              type: "is-success",
              position: "is-top",
              duration: 2000,
            });
          } else {
            Toast.open({
              message: j.msg,
              type: "is-danger",
              position: "is-top",
              duration: 2000,
            });
          }
        });
    },
    getStatus: function() {
      let that = this;
      fetch("/api/platform/status")
        .then((r) => r.json())
        .then(function(st) {
          that.platform = st;
          that.$store.commit("setPnL", Math.trunc(st.pnl));
          that.$store.commit("setPd", Math.trunc(st.pd));
          that.$store.commit("setPl", Math.trunc(st.pl));
        });
    },

    showModal: function(id) {
      this.$refs.modal.showModal({ id: id });
    },

    startGetStatus: function() {
      this.getStatus();
      let that = this;
      if (that.refreshInterval == 0) {
        return;
      }
      // console.log("platform widget refreshInterval=", that.refreshInterval);
      if (this.timer !== undefined) {
        clearTimeout(this.timer);
      }
      this.timer = setTimeout(that.startGetStatus, that.refreshInterval * 1000);
    },
  },

  created: function() {
    this.refreshInterval = this.$store.getters.refreshInterval;
    this.startGetStatus();
  },
  mounted: function() {
    let that = this;
    this.pooling = setInterval(function() {
      that.refreshInterval = that.$store.getters.refreshInterval;
    }, 1000);
  },

  beforeDestroy() {
    this.startGetStatus(0);
    clearInterval(this.pooling);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.cell {
  width: 187px;
  height: 86px;
  float: left;
  border-radius: 10px;

  margin-right: 2px;
  margin-bottom: 1px;
  position: relative;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.7);
}

.connector {
  font-size: 10px;
  text-align: left;
  font-weight: bold;
  padding-left: 2.4em;
  text-indent: -2em;
}

.widget-footer {
  position: absolute;
  bottom: 1px;
  width: 100%;
  text-align: center;
}

.log-button {
  border: 1px solid;
  border-radius: 10px;
  border-color: black;
  padding-left: 8px;
  padding-right: 8px;
  cursor: pointer;
  margin-right: 4px;
  position: absolute;
  right: 2px;
  bottom: 3px;
  float: right;
}
</style>
