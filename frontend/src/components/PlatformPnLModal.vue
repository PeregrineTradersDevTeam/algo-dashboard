<template>
  <responsive-modal
    :isActive.sync="isActive"
    :title="title"
    widthLimit="80%"
    @close="onClose"
    @resize="onResize"
  >
    <apexchart
      ref="chart"
      type="area"
      :height="chartHeight"
      :width="chartWidth"
      :options="chartOptions"
      :series="series"
    />
  </responsive-modal>
</template>

<script>
import ResponsiveModal from "@/components/ResponsiveModal.vue";
import VueApexCharts from "vue-apexcharts";

function formatTime(date) {
  var d = new Date(date * 1000),
    h = d.getHours(),
    m = d.getMinutes(),
    s = d.getSeconds();

  if (h < 10) h = "0" + h;
  if (m < 10) m = "0" + m;
  if (s < 10) s = "0" + s;

  return h + ":" + m + ":" + s;
}

export default {
  name: "PlatformPnLModal",
  components: {
    "responsive-modal": ResponsiveModal,
    apexchart: VueApexCharts,
  },
  data: function() {
    return {
      title: "",
      code: "",
      isActive: false,
      series: [],
      polling: null,
      chartRefreshRequest: 0,
      chartHeight: "80%",
      chartWidth: "100%",
      colors: [
        ["#f00", "#ed5269"],
        ["#008000", "#0d730d"],
        ["#777", "#a2a2a2"],
      ],
      chartOptions: {
        theme: { mode: "dark" },
        grid: {
          borderColor: "transparent",
          row: {
            colors: ["#131213", "#00000"], // takes an array which will be repeated on columns
            opacity: 1,
          },
        },

        stroke: {
          show: true,
          curve: "smooth",
          lineCap: "butt",
          colors: ["#ed5269"],
          width: 1,
          dashArray: 0,
        },

        fill: {
          type: "solid",
          opacity: 0.5,
          colors: ["#ff0000"],
        },
        chart: {
          type: "area",
          stacked: false,
          zoom: {
            type: "x",
            enabled: true,
            autoScaleYaxis: true,
          },
          toolbar: {
            show: true,
            autoSelected: "zoom",
            reset: false,
          },
        },
        dataLabels: {
          enabled: false,
        },
        markers: {
          size: 0,
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
          tooltip: {
            enabled: false,
          },
          type: "datetime",
          labels: {
            formatter: function(val) {
              return formatTime(val);
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
              return formatTime(val);
            },
          },
        },
      },
    };
  },
  methods: {
    showModal: function(title, code) {
      this.title = title;
      this.code = code;
      this.isActive = true;
      this.series = [{ data: [], name: "PnL" }];

      this.getCurveData();

      let that = this;
      this.polling = setInterval(function() {
        that.getCurveData();
      }, that.$store.getters.refreshInterval * 1000);
    },
    onClose: function(event) {
      clearInterval(this.polling);
      this.isActive = false;
      this.series = [];
    },

    onResize: function(screen) {
      // if (screen.type == "mobile" && screen.mode == "portrait") {
      //   this.chartHeight -= 10;
      // }
      // this.$refs.chart.render();
      //  else {
      //   this.chartHeight = "20%";
      // }
      console.log("event resize");
      //  this.$refs.chart.updateOptions();
    },

    getCurveData: function() {
      let that = this;
      fetch(
        "/api/charter/" +
          that.code +
          "/array?from=" +
          that.series[0].data.length
      )
        .then((r) => r.json())
        .then((res) => {
          if (res == null || res.data == null) {
            return;
          }
          that.series[0].data.push(...res.data);
          if (res.data.length > 0) {
            let lv = res.data[res.data.length - 1][1];
            let colorIndex = 2;
            if (lv > 0) {
              colorIndex = 1;
            }
            if (lv < 0) {
              colorIndex = 0;
            }
            that.chartOptions.fill.colors = [that.colors[colorIndex][0]];
            that.chartOptions.stroke.colors = [that.colors[colorIndex][1]];
            that.$refs.chart.updateOptions(
              that.chartOptions,
              false,
              false,
              true
            );
            that.$refs.chart.updateSeries(that.series, false);
          }
        });
    },
  },
  created: function() {},
};
</script>

<style scoped></style>
