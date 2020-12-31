<template>
  <div class="widget">
    <div v-if="hasCloseListener" class="close-button" @click="$emit('close')">
      X
    </div>
    <img :src="chartURL" :width="width" :height="height" @click="showLarge" />
    <div class="title">
      <p>{{ title }}</p>
    </div>
    <div class="depth-group">
      <div class="depth">&#8734;</div>
      <div class="depth">1h</div>
    </div>
    <PnLChartModal ref="modal" />
  </div>
</template>

<script>
import PnLChartModal from "@/components/PnLChartModal.vue";

export default {
  name: "PnLChartWidget",
  components: {
    PnLChartModal,
  },
  props: {
    title: {
      type: String,
      default: "PnL",
    },
    codes: String,
    height: Number,
    width: Number,
    renderSeq: Number,
  },
  data: function() {
    return {
      depth: 0, // all items
    };
  },

  methods: {
    showLarge: function() {
      this.$refs.modal.showModal(
        "Profit & Loss: " + this.title,
        this.codes,
        ""
      );
    },
  },
  computed: {
    hasCloseListener: function() {
      return this.$listeners && this.$listeners.close;
    },
    chartURL: function() {
      return (
        "/api/charter/" +
        this.codes +
        "/img?w=" +
        this.width +
        "&h=" +
        this.height +
        "&depth=" +
        this.depth +
        "&r=" +
        this.renderSeq
      );
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.title {
  position: absolute;
  width: 100%;
  top: 4px;
  /* z-index: 1px; */
}

.title p {
  font-size: 8pt;
  color: yellow;
  text-align: center;
}

.widget {
  width: 186vpx;
  height: 84px;
  float: left;

  margin-right: 2px;
  margin-bottom: 1px;

  text-align: left;
  display: block;
  border: 1px solid grey;
  position: relative;
}

.depth-group {
  position: absolute;
  display: inline;
  /* z-index: 1; */
  top: 67px;
  left: 0px;
  display: block;
  font-size: 7pt;
}

.depth {
  width: 15px;
  height: 15px;
  background-color: rgba(0, 0, 60, 0.3);
  opacity: 30%;
  text-align: center;
  margin-left: 4px;
  display: block;
  float: left;
}

.close-button {
  /* z-index: 20; */
  position: absolute;
  font-size: 12px;
  top: 0px;
  right: 0px;
  width: 14px;
  height: 17px;
  text-align: center;
  vertical-align: middle;
  background-color: transparent;
  cursor: pointer;
  color: white;
  opacity: 1;
}

.close-button:hover {
  color: white;
  border-radius: 9px;
  -webkit-transition-duration: 0.4s;
  /* Safari */
  transition-duration: 0.4s;
  background-color: blue;
  opacity: 1;
}
</style>
