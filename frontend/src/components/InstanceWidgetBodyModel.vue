<template>
  <div>
    <div class="cell top tooltip" tooltip="Position Twap">
      {{ m.pt }}
    </div>
    <div class="cell left top tooltip" tooltip="Position Hedge">
      {{ m.ph }}
    </div>
    <div class="cell left top tooltip last" tooltip="Position Option">
      {{ m.po }}
    </div>
    <div class="cell top tooltip" tooltip="Bid">
      {{ m.b }}
    </div>
    <div class="cell top left tooltip" tooltip="Ask">
      {{ m.a }}
    </div>
    <div class="cell top left tooltip last" tooltip="Last">
      {{ m.l }}
    </div>
    <div class="cell top tooltip" tooltip="Option Bid">
      {{ m.ob }}
    </div>
    <div class="cell top left tooltip" tooltip="Option Ask">
      {{ m.oa }}
    </div>
    <div class="cell top left tooltip last" tooltip="Share Average Hedge Price">
      {{ m.sahp }}
    </div>
    <div class="cell top bottom tooltip" tooltip="Target Quantity Per Period">
      {{ m.tqpp }}
    </div>
    <div
      class="cell top bottom left tooltip"
      tooltip="Target Position Flow Twap"
    >
      {{ m.tpft }}
    </div>
    <div class="cell top bottom left tooltip last" tooltip="Remaining Shots">
      {{ m.rs }}
    </div>
    <div class="cell chart" @click="$emit('open-chart')">
      <font-awesome-icon icon="chart-area" />
    </div>
  </div>
</template>

<script>
export default {
  name: "InstanceWidgetBody",
  components: {},
  props: {
    m: Object,
  },
  methods: {
    showTooltip(event) {
      console.log(event);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.b {
  display: flex;
}
.cell {
  width: 60px;
  height: 13px;
  float: left;
  text-align: right;
  padding-right: 1px;
  line-height: 14px;
  color: rgba(0, 0, 0, 0.7);
}

.top {
  border-top: 1px lightgray solid;
}

.left {
  border-left: 1px lightgray solid;
}
.last {
  width: 62px;
}

.bottom {
  border-bottom: 1px lightgray solid;
}

.cell:hover {
  background-color: blueviolet;
  color: white;
}

.chart {
  text-align: center;
  opacity: 0.3;
  width: 62px;
  height: 14px;
  float: left;
  padding-right: 1px;
}
.chart:hover {
  opacity: 1;
  cursor: pointer;
  background-color: blue;
  color: white;
}

.tooltip:hover:before {
  content: attr(tooltip);
  position: absolute;
  padding: 5px 10px;
  margin: -29px 0 0 -62px;
  background: blueviolet;
  color: white;
  border-radius: 3px;
}

.tooltip:hover:after {
  content: "";
  position: absolute;
  margin: -9px 0 0 -60px;
  width: 0;
  height: 0;
  border-left: 10px solid transparent;
  border-right: 10px solid transparent;
  border-top: 10px solid blueviolet;
}
.tag {
  cursor: pointer;
}
</style>
