<template>
  <responsive-modal
    name="window-instance-chart"
    :isActive.sync="isActive"
    :title="title"
    widthLimit="80%"
    @close="clean"
    @show="fixChartSizes"
  >
    <template v-slot:header>
      <button class="button is-small " @click="deselect()">
        Show All
      </button>
      <button class="button is-small" disabled @click="deselect()">
        Focus
      </button>
    </template>

    <div style="display:block; height:61%;">
      <apexchart
        ref="lastChart"
        type="line"
        height="100%"
        :options="lastOptions"
        :series="lastSeries"
        style="margin-top:-6px; "
      />
    </div>
    <div style="display:block; height:3px; background-color:transparent;"></div>
    <div style="display:block; height:30%">
      <apexchart
        ref="pnlChart"
        type="area"
        height="100%"
        :options="pnlOptions"
        :series="pnlSeries"
      />
    </div>
  </responsive-modal>
</template>

<script>
import VueApexCharts from "vue-apexcharts";
import ResponsiveModal from "@/components/ResponsiveModal.vue";

const minIndex = 0;
const maxIndex = 1;

function formatDate(date) {
  let d = new Date(date * 1000);
  let h = d.getHours(),
    m = d.getMinutes(),
    s = d.getSeconds();

  if (h < 10) h = "0" + h;
  if (m < 10) m = "0" + m;
  if (s < 10) s = "0" + s;

  return h + ":" + m + ":" + s;
}
export default {
  name: "InstanceChart",
  components: {
    apexchart: VueApexCharts,
    "responsive-modal": ResponsiveModal,
  },
  data: function() {
    return {
      annotationTemplate: {
        points: [
          {
            x: 1602371333,
            y: -2.5,
            marker: {
              size: 3,
              fillColor: "#fff",
              strokeColor: "red",
              radius: 3,
            },
            // label: {
            //   //borderColor: "#FF4560",
            //   borderColor: "",
            //   borderWidth: 0,
            //   offsetY: 26,
            //   offsetX: 26,
            //   style: {
            //     color: "#f00",
            //     //background: "#FF4560",
            //   },

            //   text: "-2.500",
            //   textAnchor: "middle",
            //   padding: {
            //     top: 60,
            //   },
            // },
          },
          {
            x: 1602371324,
            y: 13,
            marker: {
              size: 3,
              fillColor: "#fff",
              strokeColor: "#008000",
              radius: 3,
            },
            // label: {
            //   //borderColor: "#FF4560",
            //   borderWidth: 0,
            //   offsetY: 1,
            //   offsetX: 20,
            //   style: {
            //     color: "#008000",
            //     //background: "#FF4560",
            //   },

            //   text: "13.000",
            //   textAnchor: "middle",
            //   padding: {
            //     top: 60,
            //   },
            // },
          },
        ],

        yaxis: [
          {
            y: -2.5,
            borderColor: "#FF0000",
            strokeDashArray: 5,
            label: {
              borderColor: "#00FF00",
              borderWidth: 0,
              offsetY: 25,

              style: {
                color: "#f00",
                //background: "#00E396",
              },
              text: "-2.500",
            },
          },
          {
            y: 13,
            borderColor: "#00FF00",
            strokeDashArray: 5,
            label: {
              borderColor: "#00E396",
              borderWidth: 0,
              offsetY: -5,
              style: {
                color: "#008000",
                //background: "#00E396",
              },
              text: "13.000",
            },
          },
        ],
      },

      lastOptions: {
        theme: { mode: "dark" },
        stroke: {
          show: true,
          curve: ["smooth", "stepline", "stepline", "stepline"],
          lineCap: "butt",
          width: 0.75,
          dashArray: [0, 5, 0, 0],
        },
        colors: [],
        legend: {
          show: false,
        },
        chart: {
          id: "chartlast",
          type: "line",
          toolbar: {
            show: false,
            autoSelected: "pan",
          },
        },
        dataLabels: {
          enabled: false,
        },
        markers: {
          size: 0,
        },
        grid: {
          row: {
            //  colors: ["#f3f3f3", "transparent"], // takes an array which will be repeated on columns
            colors: ["#131213", "#00000"], // #2d2d2dtakes an array which will be repeated on columns
            opacity: 1,
          },
        },

        yaxis: {
          labels: {
            formatter: function(val) {
              return val.toFixed(3);
            },
          },
          title: {
            text: this.yTitle,
          },
        },
        xaxis: {
          type: "datetime",
          labels: {
            formatter: function(val) {
              return formatDate(val);
            },
          },
          tooltip: {
            enabled: false,
          },
        },
        tooltip: {
          custom: null,

          shared: true,
          y: {
            formatter: function(val) {
              return val.toFixed(3);
            },
          },
          x: {
            formatter: function(val) {
              return formatDate(val);
            },
          },
        },
      },
      pnlOptions: {
        theme: { mode: "dark" },
        stroke: {
          show: true,
          curve: "smooth",
          width: 0.75,
        },
        fill: {
          type: "solid",
          opacity: 1,
          colors: ["#FF0000"],
        },

        annotations: {},

        colors: ["silver"],
        chart: {
          type: "area",
          toolbar: {
            show: false,
            autoSelected: "selection",
          },
          brush: {
            enabled: true,
            target: "chartlast",
          },
          selection: {
            enabled: true,
            xaxis: {
              min: null,
              max: null,
            },
            fill: {
              color: "#fff",
              opacity: 0.25,
            },
          },
        },
        grid: {
          borderColor: "transparent",
          row: {
            colors: ["#131213", "#00000"], // #2d2d2dtakes an array which will be repeated on columns
            opacity: 1,
          },
        },

        dataLabels: {
          enabled: false,
        },
        markers: {},

        yaxis: {
          labels: {
            formatter: function(val) {
              return val.toFixed(3);
            },
          },
          title: {
            text: this.yTitle,
          },
        },
        xaxis: {
          type: "datetime",
          labels: {
            formatter: function(val) {
              return formatDate(val);
            },
          },
        },
        tooltip: {
          shared: false,
          y: {
            formatter: function(val) {
              return val.toFixed(3);
            },
          },
          x: {
            formatter: function(val) {
              return formatDate(val);
            },
          },
        },
      },
      isActive: false,

      // complient with Chart requirements
      lastSeries: [],
      pnlSeries: [],

      title: "",
      width: 10,
      height: 10,
      code: "",
      curves: [],
      polling: null,
      isMobile: false,
      pnlChartHeight: 200,
      orientationAngle: 0,
    };
  },
  methods: {
    fixChartSizes: function(e) {},
    tooltipFunc: function({ series, seriesIndex, dataPointIndex, w }) {
      // console.log(series);
      // console.log(seriesIndex);
      // console.log(dataPointIndex);
      //console.log(w);
      let x = this.lastSeries[seriesIndex].data[dataPointIndex][0];

      let res = `<div class="apexcharts-theme-light">
        <div class="apexcharts-tooltip-title">
          ${formatDate(x)}
        </div>`;

      for (let i = 0; i < w.globals.seriesNames.length; i++) {
        let name = w.globals.seriesNames[i];
        let color = w.config.colors[i];
        let val = "-";
        let arr = undefined;
        if (name == "pnl") {
          arr = this.pnlSeries[0].data;
        } else {
          arr = this.lastSeries[i].data;
        }

        if (arr.length == 0) {
          continue;
        }

        if (i != seriesIndex) {
          let j = 0;
          for (; j < arr.length; j++) {
            if (arr[j][0] > x) {
              // one step back to get previous value
              j -= j > 0 ? 1 : 0;
              break;
            }
          }
          if (j == arr.length) {
            j--;
          }
          val = arr[j][1].toFixed(3);
        } else {
          val = series[seriesIndex][dataPointIndex].toFixed(3);
        }

        res += `
        <div class="apexcharts-tooltip-series-group apexcharts-active" style="order: 1; display: flex;">
          <span class="apexcharts-tooltip-marker" style="background-color: ${color}"></span>
          <div class="apexcharts-tooltip-text" style="font-family: Helvetica, Arial, sans-serif; font-size: 12px;">
            <div class="apexcharts-tooltip-y-group">
              <span class="apexcharts-tooltip-text-label">${name}: </span>
              <span class="apexcharts-tooltip-text-value">${val}</span>
            </div>
            <div class="apexcharts-tooltip-z-group">
              <span class="apexcharts-tooltip-text-z-label"></span>
              <span class="apexcharts-tooltip-text-z-value"></span>
            </div>
          </div>
        </div>`;
      }

      res += "</div> ";
      return res;
    },

    updatePnlOptions: function() {
      this.$refs.pnlChart.updateOptions(this.pnlOptions, false, false, true);
    },

    focusSelection: function() {},
    deselect: function() {
      if (this.isMobile) {
        this.selectAll();
        return;
      }
      this.pnlOptions.chart.selection.xaxis.min = null;
      this.pnlOptions.chart.selection.xaxis.max = null;
      this.updatePnlOptions();
    },
    selectAll: function() {
      this.pnlOptions.chart.selection.xaxis.min = this.pnlSeries[0].data[0][0];
      this.pnlOptions.chart.selection.xaxis.max = this.pnlSeries[0].data[
        this.pnlSeries[0].data.length - 1
      ][0];

      this.updatePnlOptions();
    },
    showModal: function(shortcode, code, lastCurves) {
      if (window.orientation == 90 || window.orientation == -90) {
        this.isMobile = true;
        this.pnlChartHeight = 150;
      } else {
        this.isMobile = window.matchMedia(
          "only screen and (max-width: 492px)"
        ).matches;
        this.pnlChartHeight = 200;
      }

      this.title = shortcode;
      this.code = code;
      this.curves = lastCurves.split(",");
      this.lastSeries = [];
      this.lastOptions.colors = [];
      this.lastOptions.tooltip.custom = this.tooltipFunc;
      this.pnlSeries = [{ data: [], name: "pnl" }];
      for (let i = 0; i < this.curves.length; i++) {
        this.lastSeries.push({ data: [], name: this.curves[i] });
        this.lastOptions.colors.push("#7BD39A");
      }
      this.width = window.screen.availWidth - 80;
      this.isActive = true;
      this.refreshSeries();
    },
    clean: function() {
      clearInterval(this.polling);
      this.lastSeries = [];
      this.pnlSeries = [];
    },

    assignPnLAnnotations: function() {
      if (!this.pnlSeries[0] || this.pnlSeries[0].data < 2) {
        this.pnlOptions[0].annotations = {};
        return;
      }
      //console.log("pnl:", this.pnlSeries[0].data);
      let min = 0;
      let max = 0;
      for (let i = 0; i < this.pnlSeries[0].data.length; i++) {
        if (this.pnlSeries[0].data[i][1] > this.pnlSeries[0].data[max][1]) {
          max = i;
          continue;
        }
        if (this.pnlSeries[0].data[i][1] < this.pnlSeries[0].data[min][1]) {
          min = i;
        }
      }
      let a = this.annotationTemplate;
      //console.log("min-max", min, max);
      if (max == min) {
        this.pnlOptions[0].annotations = {};
        return;
      }
      // console.log("a", a);
      // console.log(a.points);
      // this.annotationTemplate.points[0].x = 11;
      // this.annotationTemplate.points[0].y = 12;

      //this.pnlOptions[0].annotations = annotationTemplate;
      a.points[minIndex].x = this.pnlSeries[0].data[min][0];
      a.points[minIndex].y = this.pnlSeries[0].data[min][1];
      //a.points[minIndex].label.text = this.pnlSeries[0].data[min][1].toFixed(3);

      a.points[maxIndex].x = this.pnlSeries[0].data[max][0];
      a.points[maxIndex].y = this.pnlSeries[0].data[max][1];
      //a.points[maxIndex].label.text = this.pnlSeries[0].data[max][1].toFixed(3);

      a.yaxis[minIndex].y = this.pnlSeries[0].data[min][1];
      a.yaxis[minIndex].label.text = this.pnlSeries[0].data[min][1].toFixed(3);

      a.yaxis[maxIndex].y = this.pnlSeries[0].data[max][1];
      a.yaxis[maxIndex].label.text = this.pnlSeries[0].data[max][1].toFixed(3);

      this.pnlOptions.annotations = a;
    },

    refreshSeries: function() {
      let url =
        "/api/charter/" + this.code + "/multiarray?a=" + this.curves.join(",");

      let that = this;
      fetch(url)
        .then((r) => r.json())
        .then((res) => {
          if (res == null) {
            return;
          }
          for (let k in res) {
            if (k == "pnl") {
              if (res.pnl.data.length == 0) {
                that.pnlOptions.chart.brush.enabled = false;
                that.lastOptions.chart.continue;
              }
              that.pnlOptions.chart.brush.enabled = true;
              that.pnlSeries[0].data.push(...res.pnl.data);
              that.pnlSeries[0].name = res.name;
              that.assignPnLAnnotations();
              that.$refs.pnlChart.updateSeries(that.pnlSeries, false);
              if (that.isMobile) {
                that.selectAll();
              }
              that.pnlOptions.fill.colors[0] =
                res.pnl.data[res.pnl.data.length - 1][1] < 0
                  ? "rgba(255,  0, 0, 0.3)"
                  : "rgba(0,  255, 0, 0.3)";

              that.pnlOptions.colors[0] = that.pnlOptions.fill.colors[0];
              that.updatePnlOptions();
              continue;
            }

            const idx = that.curves.indexOf(k);
            that.lastSeries[idx].data.push(...res[k].data);
            that.lastSeries[idx].data.name = res[k].name;
            that.lastOptions.colors[idx] = res[k].color;
            if (k == "op") {
              that.lastOptions.dashArray = 5;
            }
          }
          that.$refs.lastChart.updateSeries(that.lastSeries, false);
          console.log(that.lastOptions);
        });
    },
  },
};
</script>

<style scoped>
.toolbar {
  height: 50px;
  background-color: whitesmoke;
  width: 100%;
}

.tooltip {
  border-radius: 5px;
  box-shadow: 2px 2px 6px -4px #999;
  cursor: default;
  font-size: 14px;
  left: 62px;
  opacity: 0;
  pointer-events: none;
  position: absolute;
  top: 20px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  white-space: nowrap;
  /*z-index: 12;*/
  transition: 0.15s ease all;
}
</style>
