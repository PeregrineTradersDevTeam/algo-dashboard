<template>
  <div>
    <div
      v-if="selected == ''"
      class="cell"
      v-bind:style="{ border: borderStyle }"
    >
      <div
        class="group"
        v-for="item in groups"
        v-bind:key="item.name"
        @click="selected = item.name"
      >
        <span class="l">{{ item.name }}</span>
        <span
          v-bind:class="{
            'value value-profit': item.val > 0,
            value: item.val > 0,
            'value value-loss': item.val < 0,
            value: item.val == 0,
          }"
          >{{ item.val }}</span
        >
      </div>
    </div>
    <div v-else>
      <PnLChartWidget
        :title="selected"
        :codes="'PNL:BU:' + selected"
        :renderSeq="renderSeq"
        :width="185"
        :height="84"
        @close="selected = ''"
      />
    </div>
  </div>
</template>

<script>
import PnLChartWidget from "@/components/PnLChartWidget.vue";

export const LAST = 0;
export const MIN = 1;
export const MAX = 2;

const urlvar = ["", "min", "max"];

export default {
  name: "GroupPnLWidget",
  props: {
    mode: {
      default: 0,
      type: Number,
    },
    hideChart: Boolean,
  },
  components: {
    PnLChartWidget,
  },

  data: function() {
    return {
      renderSeq: 0,
      polling: null,
      groups: [],
      selected: "",
    };
  },
  watch: {
    mode: function() {
      this.getData();
    },
    hideChart: function() {
      this.selected = "";
    },
  },
  computed: {
    borderStyle: function() {
      let res = "1px solid ";
      switch (this.mode) {
        case MIN:
          res += "red";
          break;
        case MAX:
          res += "green";
          break;
        default:
          res += "silver";
          break;
      }
      return res;
    },
  },

  methods: {
    getData: function() {
      this.renderSeq++;
      let that = this;
      fetch(
        "/api/portfolio/group-pnl?except=PNL:BUOYANCY&type=" + urlvar[that.mode]
      )
        .then((r) => r.json())
        .then(function(items) {
          if (items.length == that.groups.length) {
            for (let i = 0; i < items.length; i++) {
              that.$set(that.groups, i, items[i]);
            }
            return;
          }

          const diff = items.length - that.groups.length;
          if (diff < 0) {
            // read less then we have already
            for (let i = 0; i < items.length; i++) {
              that.$set(that.groups, i, items[i]);
            }
            that.groups.splice(-1 * diff, that.groups.length);
            return;
          }

          // diff > 0
          // read less then we have already
          let i = 0;
          for (; i < that.groups.length; i++) {
            that.$set(that.groups, i, items[i]);
          }
          for (; i < items.length; i++) {
            that.groups.push(items[i]);
          }
          return;
        });
    },
  },
  mounted: function() {
    this.getData();
    let that = this;
    this.pooling = setInterval(function() {
      that.getData();
    }, 1000);
  },

  beforeDestroy() {
    clearInterval(this.pooling);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.cell {
  width: 185px;
  height: 84px;
  float: left;

  margin-right: 2px;

  margin-bottom: 1px;
  position: relative;
  font-size: 12px;
}

.group {
  float: left;
  width: 33%;
  font-size: 10px;
  text-align: left;
  font-weight: bold;
}

.group:nth-child(odd) {
  background: #2c2c2c;
}
.group:nth-child(3n) {
  width: 34%;
}

.group:hover {
  background-color: lightskyblue;
  cursor: pointer;
}

.l {
  font-size: 10px;
  color: gray;
  padding-left: 2px;
}

.value {
  font-size: 10px;
  text-align: right;
  float: right;
  padding-right: 5px;
  color: rgba(255, 255, 255, 0.2);
}

.value:nth-child(3n) {
  font-size: 10px;
  text-align: right;
  margin-right: 0px;
}

.value-profit {
  color: #11e95c;
}

.value-loss {
  color: red;
}
</style>
