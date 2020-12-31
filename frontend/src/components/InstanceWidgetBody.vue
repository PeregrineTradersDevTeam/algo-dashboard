<template>
  <div v-bind:style="{ 'background-color': bgColor }">
    <InstanceWidgetBodyRisk v-if="mode == 'risk'" :data="m" />
    <InstanceWidgetBodyChart
      v-else-if="mode == 'chart'"
      :id="m.id"
      @open-chart="$emit('open-chart')"
    >
    </InstanceWidgetBodyChart>
    <InstanceWidgetBodyModel v-else @open-chart="$emit('open-chart')" :m="m" />
  </div>
</template>

<script>
import InstanceWidgetBodyChart from "@/components/InstanceWidgetBodyChart.vue";
import InstanceWidgetBodyRisk from "@/components/InstanceWidgetBodyRisk.vue";
import InstanceWidgetBodyModel from "@/components/InstanceWidgetBodyModel.vue";

export default {
  name: "InstanceWidgetBody",
  components: {
    InstanceWidgetBodyChart,
    InstanceWidgetBodyRisk,
    InstanceWidgetBodyModel,
  },
  props: {
    m: Object,
    mode: String,
  },
  data: function() {
    return {
      strategyBackgroundColor: {
        BUOYANCY: "silver",
        TREND_TRADER: "#E8CBED",
        SCRIP: "",
        PAIR_TRADER: "",
      },
    };
  },
  computed: {
    bgColor() {
      return this.strategyBackgroundColor[this.m.atype];
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
