<template>
  <div class="page-header">
    <div class="tags">
      <Tag short label="zone" :value="zone" valueBC="#3298dc" />
      <Tag short label="site" :value="site" valueBC="#48c774" />
      <Tag short label="env" :value="env" valueBC="#48c774" />
      <Tag short label="platform" :value="platform" valueBC="#7957d5" />
      <Tag
        short
        label="target"
        :value="ppnl"
        valueBC="#ffdd57"
        valueColor="black"
      />
    </div>
    <div class="status">
      <PnlBadge />
      <InstanceStatusCounter
        location="header"
        v-if="screen.mode != 'portrait' && screen.width > 578"
      />
    </div>

    <div class="button-group">
      <ResponsiveButton
        title="Launcher"
        @click="launch"
        :isShort="screen.width < 688"
      />
      <ResponsiveButton
        title="Update"
        @click="update"
        :isShort="screen.width < 688"
      />
      <ResponsiveButton
        title="Close ALL"
        @click="confirm"
        bgColor="#f14668"
        color="white"
        :isShort="screen.width < 688"
      />
      <LauncherModal ref="modal" />
    </div>
  </div>
</template>

<script>
import InstanceStatusCounter from "@/components/InstanceStatusCounter.vue";
import Tag from "@/components/Tag.vue";
import PnlBadge from "@/components/PnlBadge.vue";
import ResponsiveButton from "@/components/ResponsiveButton.vue";
import LauncherModal from "@/components/LauncherModal.vue";

export default {
  name: "TheHeader",
  components: {
    InstanceStatusCounter,
    Tag,
    PnlBadge,
    ResponsiveButton,
    LauncherModal,
  },
  data: function() {
    return {
      isLauncherActive: false,
      isUpdaterActive: false,
      zone: "",
      site: "",
      platform: "",
      env: "",
      ppnl: "",
    };
  },
  computed: {
    screen() {
      return this.$store.getters.sp;
    },
    lastPnL() {
      return this.$store.state.lastPnL;
    },
    lastPd() {
      return this.$store.state.lastPd;
    },
    lastPl() {
      return this.$store.state.lastPl;
    },
  },

  methods: {
    launch() {
      this.$refs.modal.showModal(
        "Algo Launcher",
        "Start ALL",
        "launch",
        this.$store.getters.holidayToday
      );
    },
    update() {
      this.$refs.modal.showModal(
        "Algo Updater",
        "Send Update",
        "update",
        this.$store.getters.holidayToday
      );
    },

    confirm() {
      this.$confirm({
        title: "Platform",
        message: "Close all instances?",
        button: { yes: "Yes", no: "No" },
        callback: (confirm) => {
          if (confirm) this.closeAll();
        },
      });
    },

    closeAll() {
      let that = this;
      fetch("/api/close-all", { method: "POST" })
        .then((r) => r.json())
        .then(function(j) {
          if (j.received === true) {
            that.$toast.success("CLOSE ALL received by server!");
          } else {
            that.$toast.error("CLOSE ALL has not been received by server!");
          }
        });
    },

    getCofig: function() {
      let that = this;
      fetch("/api/config")
        .then((r) => r.json())
        .then(function(j) {
          that.zone = j["zone"];
          that.platform = j["platform"];
          that.site = j["site"];
          that.env = j["env"];
          that.ppnl = j["ppnl"];
          that.$store.commit("setPlatformBuild", j["build"]);
          that.$store.commit("setPdlc", j["pdlc"]);
          document.title = that.env + ":" + that.zone + ":" + that.platform;
        });
    },
  },
  created: () => {
    document.title = "Dashboard";
  },

  mounted: function() {
    this.getCofig();
    let that = this;
    this.pooling = setInterval(function() {
      that.getCofig();
    }, 30000);
  },
  beforeDestroy() {
    clearInterval(this.pooling);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.page-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  flex-wrap: nowrap;
  flex: 1;
  padding-top: 2px;
  background-color: #1e1e1e;
  height: 25px;
}

@media (min-width: 768px) {
  .hideable {
    display: inline;
  }
}

.tags {
  display: flex;
  flex-flow: nowrap;
}

.status {
  display: flex;
  flex-flow: nowrap;
  margin-left: 2px;
}

.key-buttons {
  float: right;
}

.short-button-text {
  display: none;
}
/*
@media only screen and (max-width: 1040px) {
  .separator {
    width: 2px;
  }
}

@media only screen and (max-width: 804px) {
  .hider {
    display: none;
  }
}
.tight {
  padding-left: 6px;
  padding-right: 6px;
/* }

@media only screen and (max-width: 579px) {
  .tight {
    padding-left: 3px;
    padding-right: 3px;
  }
}

@media only screen and (max-width: 894px) {
  .short-button-text {
    display: inline-block;
  }
  .full-button-text {
    display: none;
  }
  .tight {
    padding-left: 3px;
    padding-right: 3px;
  }
} */
.button-group {
  display: flex;
  flex-wrap: nowrap;
  justify-content: space-between;
  flex-shrink: 0;
  margin: 0 4px 0px 4px;
  min-width: 80px;
}

/* .pnl-counter {
  height: 24px;
  padding: 1px 6px 1px 6px;
  font-weight: bold;
  font-size: 14px;
  border: 1px;
  border-style: solid;
  border-color: black;
  color: white;
}
.pnl-positive {
  background-color: green;
}

.pnl-negative {
  background-color: red;
} */
</style>
