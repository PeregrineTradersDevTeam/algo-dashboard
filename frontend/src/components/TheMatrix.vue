<template>
  <div class="matrix" v-touch:swipe="swipeHandler">
    <GroupPnLWidget :mode="groupMode" :hideChart="group.hideChart" />

    <PnLChartWidget
      title="Platform"
      codes="PNL:BUOYANCY"
      :renderSeq="renderSeq"
      :width="185"
      :height="84"
    />

    <PlatformWidget
      ref="modal"
      title="Platform Widget"
      :isStopped="isPlatformStopped"
    />

    <InstanceWidget
      v-for="m in matrixm"
      v-bind:key="m.id"
      :m="m"
      :i="matrixi[m.id]"
      :mode="currmode"
      :hideBody="collapsed"
      @close="closeAlgo(m.id)"
      @open="showModal(m.id)"
      @open-chart="showChartModal(m.id)"
    />
    <InstanceModal ref="modal" />
    <InstanceChart ref="chartModal" />
  </div>
</template>

<script>
import InstanceModal from "@/components/InstanceModal.vue";
import PnLChartWidget from "@/components/PnLChartWidget.vue";
import PlatformWidget from "@/components/PlatformWidget.vue";
import { ToastProgrammatic as Toast } from "buefy";
import InstanceChart from "@/components/InstanceChart.vue";

import GroupPnLWidget from "@/components/GroupPnLWidget.vue";

import InstanceWidget from "@/components/InstanceWidget.vue";

const largeModes = "l,op,lo,hi,pnl";

export const SORT_BY_ID = 0;
export const SORT_BY_PNL = 1;
export const SORT_BY_STATUS = 2;

export default {
  name: "TheMatrix",
  components: {
    InstanceModal,
    PlatformWidget,
    PnLChartWidget,
    InstanceChart,
    GroupPnLWidget,
    InstanceWidget,
  },
  props: {
    matrixSortMode: Number,
    groupMode: Number,
    collapsed: Boolean,
    currmode: String,
  },
  data: function() {
    return {
      largeCurves: "l,op,lo,hi,pnl", // the order of symbols are important!
      renderSeq: 0,
      lpmode: false,
      hover: false,
      matrixm: [],
      matrixi: {},
      portfolio: new Map(),
      group: {
        title: "lst",
        titles: ["lst", "min", "max"],
        selected: 0,
        current: 0,
        modes: ["", "min", "max"],
        mode: "",
        hideChart: true,
      },
      sort: {
        title: "ID", // current button title
        mode: "id", // current matrix sort mode
        modes: ["id", "pnl", "status"], // supported sort modes
        titles: ["ID", "PnL", "ST"], //
      },
      config: {
        strategy: {
          BUOYANCY: {},
          TREND_TRADER: {},
          SCRIP: {},
          PAIR_TRADER: {},
        },
        strategyBackgroundColor: {
          BUOYANCY: "silver",
          TREND_TRADER: "#E8CBED",
          SCRIP: "",
          PAIR_TRADER: "",
        },
        instanceStatusColor: {},
      },
      refreshInterval: 1,
      poolingM: null,
      poolingI: null,
      isPlatformStopped: true,
      failedIntances: {},
      isMinimaized: false,
      isConfig: false,
    };
  },
  computed: {
    sevcolor() {
      return "red";
    },
  },
  watch: {
    matrixSortMode(n) {
      this.doSort(this.matrixm, n);
    },
  },
  methods: {
    showModal(id) {
      this.$refs.modal.showModal({ id: id });
    },
    showChartModal(id) {
      this.$refs.chartModal.showModal(id, "M:" + id, this.largeCurves);
    },

    showLarge: function(shortcode, code) {
      this.$refs.instanceChartModal.showModal(shortcode, code, largeModes);
    },

    // changeSort: function() {
    //   let idx = this.sort.modes.indexOf(this.sort.mode);
    //   idx++;
    //   if (idx == this.sort.modes.length) {
    //     idx = 0;
    //   }
    //   this.sort.mode = this.sort.modes[idx];
    //   this.sort.title = this.sort.titles[idx];
    //   this.getMatrixM();
    //   localStorage.setItem("matrix-sort", this.sort.mode);
    // },

    // changeGroupMode: function() {
    //   let idx = this.group.modes.indexOf(this.group.mode);
    //   idx++;
    //   if (idx == this.group.modes.length) {
    //     idx = 0;
    //   }
    //   this.group.mode = this.group.modes[idx];
    //   this.group.title = this.group.titles[idx];
    // },

    swipeHandler: function(direction) {
      this.$emit("swipe", direction);
      //this.nextMode({ left: false, right: true }[direction]);
    },

    // nextMode: function(forward) {
    //   let idx = this.matrixMode.indexOf(this.currmode);
    //   if (idx < 0) {
    //     this.currmode = this.matrixMode[0];
    //     return;
    //   }
    //   idx += forward ? 1 : -1;

    //   if (idx == this.matrixMode.length) {
    //     idx = 0;
    //   } else if (idx < 0) {
    //     idx = this.matrixMode.length - 1;
    //   }
    //   this.currmode = this.matrixMode[idx];
    // },

    getMatrixI: function() {
      let that = this;
      fetch("/api/matrix/i")
        .then((r) => r.json())
        .then(function(items) {
          if (items) {
            that.matrixi = {};
            that.matrixi = items;
          }
        });
    },

    getIntanceAttr: function(id, attr) {
      let i = this.matrixi[id];
      if (i) {
        if (i[attr]) {
          return i[attr];
        }
      }
      return "";
    },

    doSort(items, mode) {
      // if (this.currmode == "risk") {
      //   return;
      // }
      switch (mode) {
        case SORT_BY_ID:
          items.sort((a, b) => (a.id < b.id ? -1 : a.id > b.id ? 1 : 0));
          break;
        case SORT_BY_PNL:
          items.sort((a, b) => {
            if (a.pnl == 0 && b.pnl == 0) {
              return a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
            }
            if (a.pnl == 0) {
              return 1;
            }
            if (b.pnl == 0) {
              return -1;
            }
            return a.pnl < b.pnl ? -1 : 1;
          });
          break;
        case SORT_BY_STATUS:
          items.sort((a, b) => {
            if (a.status == b.status) {
              return a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
            }

            return a.status < b.status ? -1 : a.status > b.status ? 1 : 0;
          });
          break;
      }
    },
    applyMatrixM(items) {
      if (this.currmode == "risk") {
        // just replace values without sorting.
        // and replacing values in the same widget positions.
        for (let i = 0; i < this.matrixm.length; i++) {
          for (let j = 0; j < items.length; j++) {
            if (this.matrixm[i].id == items[j].id) {
              this.$set(this.matrixm, i, items[j]);
              break;
            }
          }
        }
        return;
      }

      if (items.length == this.matrixm.length) {
        for (let i = 0; i < items.length; i++) {
          this.$set(this.matrixm, i, items[i]);
        }
        return;
      }

      const diff = items.length - this.matrixm.length;
      if (diff < 0) {
        // read less then we have already
        for (let i = 0; i < items.length; i++) {
          this.$set(this.matrixm, i, items[i]);
        }
        this.matrixm.splice(-1 * diff, this.matrixm.length);
        return;
      }

      // diff > 0
      // read less then we have already
      let i = 0;
      for (; i < this.matrixm.length; i++) {
        this.$set(this.matrixm, i, items[i]);
      }
      for (; i < items.length; i++) {
        this.matrixm.push(items[i]);
      }
    },

    getMatrixM: function() {
      let that = this;
      fetch("/api/matrix/m") // ?sort=" + that.sort.mode)
        .then((r) => r.json())
        .then(function(items) {
          let scnt = { 0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 };
          that.isPlatformStopped = true;
          for (let i = 0; i < items.length; i++) {
            if (items[i].atype == "BUOYANCY") {
              items[i]["atypeshort"] = "BU";
            }
            scnt[items[i].status]++;
            if (items[i].status < 2) {
              that.isPlatformStopped = false;
            }
            if (items[i].status == 3) {
              if (!that.failedIntances[items[i].id]) {
                that.failedIntances[items[i].id] = false;
              }
            }
          }
          if (that.currmode != "risk") {
            // sort only in model mode
            that.doSort(items, that.matrixSortMode);
          }

          that.applyMatrixM(items);

          //that.matrixm.splice(0, that.matrixm.length, ...items);

          that.$store.commit("setInstanceStatusCounter", scnt);
        });
    },

    closeAlgo: function(code) {
      let that = this;
      this.$confirm({
        title: "Instance: " + code,
        message: "Do you confirm instance close?",
        button: { yes: "Yes", no: "No" },
        callback: (confirm) => {
          if (!confirm) {
            return;
          }

          fetch(`/api/instance/${code}/stop`, { method: "POST" })
            .then((r) => r.json())
            .then(function(j) {
              let msg = `CLOSE ${code} `;
              if (j.received === true) {
                that.$toast.success(msg + "received by server!");
              } else {
                that.$toast.error(msg + "has not been received by server!");
              }
            });
        },
      });
    },
    downloadCSV: function(code) {
      window.open("/api/log/" + code + "/download");
    },
    downloadCurve: function(code) {
      window.open("/api/charter/" + code + "/download");
    },

    showModal: function(id, activeTab = 0) {
      this.$refs.modal.showModal({ id: id }, activeTab);
    },

    startRefreshMatrix: function() {
      this.renderSeq++;
      this.getMatrixM();
      let that = this;
      if (that.refreshInterval == 0) {
        return;
      }
      if (this.timer !== undefined) {
        clearTimeout(this.timer);
      }
      for (let k in this.failedIntances) {
        if (this.failedIntances[k] == false) {
          this.showModal(k, 1);
          this.failedIntances[k] = true;
          break;
        }
      }
      this.timer = setTimeout(
        that.startRefreshMatrix,
        that.refreshInterval * 1000
      );
    },
  },
  created: function() {
    this.config.instanceStatusColor = this.$store.getters.instanceStatusColor;
    this.refreshInterval = this.$store.getters.refreshInterval;
    // if (this.$store.getters.isRiskVisible === true) {
    //   this.matrixMode.push("risk");
    // }
    this.startRefreshMatrix();
  },
  mounted: function() {
    this.getMatrixI();
    this.poolingI = setInterval(function() {
      that.getMatrixI();
    }, 30000);

    // this.sort.mode = localStorage.getItem("matrix-sort");
    // let idx = this.sort.modes.indexOf(this.sort.mode);
    // if (idx >= 0) {
    //   this.sort.title = this.sort.titles[idx];
    // }

    let that = this;
    this.poolingM = setInterval(function() {
      that.refreshInterval = that.$store.getters.refreshInterval;
    }, 1000);

    window.addEventListener("keydown", (e) => {
      switch (e.code) {
        case "Escape":
          this.group.hideChart = !this.group.hideChart;
          break;
      }
    });
  },

  beforeDestroy() {
    this.startRefreshMatrix(0);
    clearInterval(this.poolingM);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.matrix {
  display: inline-block;
  width: 100%;
  position: relative;
  float: left;
  /* align-content: flex-start; */
}
/*
@media only screen and (max-width: 440px) {
  .matrix {
    width: 379px;
    display: inline-block;
  }
} */
</style>
