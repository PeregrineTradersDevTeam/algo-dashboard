<template>
  <div class="grid">
    <div tooltip="Position Delta" class="cell top tooltip">
      {{ format(data.pd) }}
    </div>
    <div tooltip="Num of Securities" class="cell left top tooltip ">
      {{ format(data.ns) }}
    </div>
    <div tooltip="Position Delta Limit" class="cell left top tooltip last">
      {{ format(data.pl) }}
    </div>
    <div tooltip="Position Delta Cash" class="cell top bottom">
      {{ format(data.pdc) }}
    </div>
    <div
      tooltip="Position Delta Limit Cash"
      class="cell top left tooltip bottom"
    >
      {{ format(data.plc) }}
    </div>
    <div class="cell top bottom left last">&nbsp;</div>
  </div>
</template>

<script>
export default {
  name: "InstanceWidgetBodyRisk",
  props: {
    data: Object,
  },
  methods: {
    format(val) {
      if (!val) {
        return "";
      }
      let n = parseFloat(val);
      if (isNaN(n)) {
        return val;
      }
      return Math.trunc(n);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.grid {
  background-color: rgba(0, 0, 0, 0.8);
  width: 100%;
  height: 70px;
}
.cell {
  width: 60px;
  height: 13px;
  float: left;
  text-align: right;
  padding-right: 1px;
  line-height: 14px;
  color: wheat;
}
.last {
  width: 62px;
}

.top {
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.left {
  border-left: 1px solid rgba(255, 255, 255, 0.05);
}

.right {
  border-right: 1px solid rgba(255, 255, 255, 0.05);
}

.bottom {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.cell:hover {
  background-color: blueviolet;
  color: white;
}

.tag {
  cursor: pointer;
}

.tooltip:hover:before {
  content: attr(tooltip);
  position: absolute;
  padding: 5px 10px;
  margin: -29px 0 0 -62px;
  background: blueviolet;
  color: white;
  border-radius: 3px;
  z-index: 3;
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
  z-index: 2;
}
</style>
