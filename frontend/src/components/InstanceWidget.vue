<template>
  <div
    class="instance-widget"
    @mouseover="hover = true"
    @mouseleave="hover = false"
    v-bind:style="{
      'background-color': config.strategyBackgroundColor[m.atype],
      height: hideBody ? 17 : 87,
    }"
  >
    <InstanceWidgetHeader
      :id="m.id"
      :pnlGroup="iGroup"
      :status="m.status"
      :currencyCode="currency"
      :pnl="m.pnl"
      :shortAlgoType="m.atypeshort"
      :hovered="hover"
      @close="$emit('close')"
      @open="$emit('open')"
    />
    <InstanceWidgetBody
      v-if="!hideBody"
      :m="m"
      :mode="mode"
      @open-chart="$emit('open-chart')"
    />
  </div>
</template>

<script>
import InstanceWidgetHeader from "@/components/InstanceWidgetHeader.vue";
import InstanceWidgetBody from "@/components/InstanceWidgetBody.vue";

export default {
  name: "InstanceWidget",
  components: {
    InstanceWidgetHeader: InstanceWidgetHeader,
    InstanceWidgetBody: InstanceWidgetBody,
  },
  props: {
    m: Object,
    i: Object,
    hideBody: Boolean,
    mode: String,
    headerColor: String,
  },
  data: function() {
    return {
      hover: false,
      config: {
        strategyBackgroundColor: {
          BUOYANCY: "silver",
          TREND_TRADER: "#E8CBED",
          SCRIP: "",
          PAIR_TRADER: "",
        },
      },
    };
  },
  computed: {
    iGroup() {
      return this.i && this.i.pnlGroup ? this.i.pnlGroup : "";
    },
    currency() {
      return this.i && this.i.currencyCode ? this.i.currencyCode : "";
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.instance-widget {
  width: 187px;
  float: left;
  margin-right: 2px;
  margin-bottom: 1px;
  white-space: nowrap;
}

.instance-widget:hover {
  -webkit-box-shadow: inset 0px 0px 4px 0px rgba(9, 13, 237, 1);
  -moz-box-shadow: inset 0px 0px 4px 0px rgba(9, 13, 237, 1);
  box-shadow: 0px 0px 10px 0px slateblue;
}

/* .separator {
  -webkit-box-shadow: 10px 0px 6px -3px rgba(9, 13, 237, 1);
  -moz-box-shadow: 10px 0px 6px -3px rgba(9, 13, 237, 1);
  box-shadow: 10px 0px 0px -3px yellow;
} */
</style>
