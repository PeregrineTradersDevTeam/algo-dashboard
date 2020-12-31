<template>
  <div class="home">
    <div class="header">
      <div class="main">
        <TheHeader />
      </div>
      <div class="risk" v-if="riskModeVisible">
        <TheHeaderRisk />
      </div>
    </div>
    <div class="content">
      <div class="matrix">
        <TheMatrix
          :groupMode="groupMode"
          :matrixSortMode="matrixSortMode"
          :viewMode="viewMode"
          :collapsed="collapseMode == 1"
          :isRiskVisible="riskModeVisible"
          :currmode="currentMatrixMode"
          @swipe="swipe"
        />
      </div>
    </div>
    <!-- <InstanceStatusCounter location="right" /> -->
    <div class="footer-f">
      <TheFooter />
    </div>

    <div class="vertical-right-panel">
      <div style="display:flex; flex-direction:row; height: 100%;">
        <InstanceStatusCounter
          location="right"
          v-if="screen.mode === 'portrait' && screen.type === 'mobile'"
        />
      </div>
      <!-- v-if="screen.width <= 578" /> -->

      <FloatingButton
        icon="calendar"
        @click="$refs.mhModal.showModal()"
        :bgColor="$store.getters.holidays.isHolidayToday ? 'red' : ''"
        :color="$store.getters.holidays.isUpdateRequired ? 'darkorange' : ''"
      />

      <div class="button-group">
        <RollingButton
          :title="['lst', 'min', 'max']"
          :selected="groupMode"
          @rolled="groupMode = $event"
          KeyCode="KeyG"
        />
        <RollingButton
          :title="['ID', 'PnL', 'ST']"
          :selected="matrixSortMode"
          localStorageKey="matrix-sort"
          @rolled="matrixSortMode = $event"
          KeyCode="KeyS"
        />

        <RollingButton
          :title="['arrow-down', 'arrow-up']"
          :selected="collapseMode"
          localStorageKey="matrix-collapsed"
          mode="icon"
          @rolled="collapseMode = $event"
          KeyCode="KeyC"
        />
        <FloatingButton
          icon="cog"
          @click="$refs.configModal.showModal($store.getters.config)"
        />
      </div>
    </div>
    <MarketHolidaysModal ref="mhModal" />
    <TheConfig ref="configModal" @save="saveConfig" />
    <!-- <LauncherModal ref="launcher" /> -->
  </div>
</template>

<script>
// @ is an alias to /src
import TheHeader from "@/components/TheHeader.vue";
import TheHeaderRisk from "@/components/TheHeaderRisk.vue";
import TheFooter from "@/components/TheFooter.vue";
import TheMatrix from "@/components/TheMatrix.vue";
import InstanceStatusCounter from "@/components/InstanceStatusCounter.vue";
import MarketHolidaysModal from "@/components/MarketHolidaysModal";
import RollingButton from "@/components/RollingButton";
import FloatingButton from "@/components/FloatingButton";
import TheConfig from "@/components/TheConfig";

const risk = "risk";
export default {
  name: "Home",
  components: {
    TheHeader,
    TheHeaderRisk,
    TheMatrix,
    TheFooter,
    InstanceStatusCounter,
    MarketHolidaysModal,
    RollingButton,
    FloatingButton,
    TheConfig,
  },
  data: function() {
    return {
      groupMode: 0,
      matrixSortMode: 0,
      viewMode: 0,
      collapseMode: 0,
      matrixModes: ["model", "chart"],
      currentMatrixMode: "model",
      isConfigActive: false,
    };
  },
  computed: {
    riskModeVisible() {
      return this.$store.state.config.isRiskVisible;
    },
    screen() {
      return this.$store.getters.sp;
    },
  },
  watch: {
    "$store.state.holidays.userNotifyRequired": function() {
      if (this.$store.state.holidays.userNotifyRequired == true) {
        this.$toast.error("Today is holiday! Check the list!", {
          timeout: 10000,
        });
      }
    },
  },
  methods: {
    saveConfig(cfg) {
      this.$store.commit("setConfig", cfg);
      if (cfg.isRiskVisible) {
        if (this.matrixModes.indexOf(risk) < 0) {
          this.matrixModes.push(risk);
        }
      } else {
        const idx = this.matrixModes.indexOf(risk);
        if (idx >= 0) {
          if (this.currentMatrixMode === this.matrixModes[idx]) {
            this.nextMode(true);
          }
          this.matrixModes.splice(idx, 1);
        }
      }
    },
    swipe(direction) {
      this.nextMode({ left: false, right: true }[direction]);
    },

    nextMode(forward) {
      let idx = this.matrixModes.indexOf(this.currentMatrixMode);
      if (idx < 0) {
        this.currentMatrixMode = this.matrixModes[0];
        return;
      }
      idx += forward ? 1 : -1;

      if (idx >= this.matrixModes.length) {
        idx = 0;
      } else if (idx < 0) {
        idx = this.matrixModes.length - 1;
      }
      this.currentMatrixMode = this.matrixModes[idx];
    },
  },
  mounted() {
    if (this.$store.getters.isRiskVisible === true) {
      this.matrixModes.push("risk");
    }

    let that = this;
    window.addEventListener(
      "keypress",
      (e) => {
        switch (e.code) {
          case "KeyV":
            that.nextMode(true);
            break;
          case "Digit1":
            that.currentMatrixMode = this.matrixModes[0];
            break;
          case "Digit2":
            that.currentMatrixMode = this.matrixModes[1];
            break;
          case "Digit3":
            if (that.matrixModes.length == 3) {
              that.currentMatrixMode = this.matrixModes[2];
            }
            break;
          default:
            break;
        }
      },
      { passive: true }
    );
  },
};
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header {
  display: flex;
  flex-direction: column;
  width: 100%;
  position: sticky;
  top: 0px;
  z-index: 1;
  background-color: #1e1e1e;
}

.main {
  flex: 0 0 28px;
}
.risk {
  height: 26px;
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  align-items: top;
}

.content {
  display: flex;
  overflow-y: hidden;
  overflow-x: hidden;
  min-height: 2em;
  width: calc(100% - 30px);
  padding-left: 1px;
}

.matrix {
  flex: none;
  width: 100%;
}
/* @media only screen and (max-width: 563px) {
  .matrix {
    width: 91%;
  }
} */

/* .right-placeholder {
  display: flex;
  position: sticky;
  width: 30px;
  z-index: 1;
  top: calc(50% - px);
  padding: 1px 1px 1px 1px;
  margin-top: 1px;
  background-color: transparent;
} */
.right {
  display: flex;
  flex-direction: column;
  flex-wrap: nowrap;
  position: sticky;
  z-index: 1;
  background-color: rebeccapurple;
}

@media only screen and (min-width: 563px) {
  .right-panel {
    display: none;
  }
}

.footer-f {
  flex-shrink: 0;
}

.pnlchart {
  width: 100vw;
  width: 100%;
}
.vertical-right-panel {
  display: block;
  position: fixed;
  right: 1px;
  top: 54px;
  width: 30px;
}

.button-group {
  display: flex;
  flex-direction: column;
  flex-wrap: nowrap;
  position: sticky;
  z-index: 1;
}
</style>
