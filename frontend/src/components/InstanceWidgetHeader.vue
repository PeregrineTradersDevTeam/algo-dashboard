<template>
  <div
    class="header"
    v-bind:style="{ 'background-color': instanceStatusColor[status] }"
    v-bind:class="{ hovered: hovered }"
  >
    <a class="product-id" v-on:click="$emit('open')">{{ productID }}</a>
    <!-- <span class="algo-type">{{ shortAlgoType ? shortAlgoType : "?" }}</span> -->
    <span style="display:inline-block;">
      <Dropdown
        :items="menu"
        @click="downloadCSV(id)"
        :title="shortAlgoType ? shortAlgoType : '?'"
      />
    </span>

    <a class="close-button" v-on:click="$emit('close')">X</a>
    <span v-bind:class="['pnl', pnlClass]">
      {{ pnl ? pnl.toFixed(2) : "" }}
    </span>
    <span class="currency">
      {{ currencyCode }}
    </span>
  </div>
</template>

<script>
import Dropdown from "@/components/Dropdown.vue";

export default {
  name: "InstanceWidgetHeader",
  components: {
    Dropdown,
  },
  props: {
    id: {
      type: String,
      required: true,
    },
    status: Number,
    pnlGroup: String,
    currencyCode: String,
    pnl: Number,
    shortAlgoType: String,
    hovered: Boolean,
  },
  data: function() {
    return {
      instanceStatusColor: {},
      menu: ["Download last CSV", "Download previous CSVs"],
    };
  },
  methods: {
    downloadCSV: function(code) {
      window.open("/api/log/" + code + "/download");
    },
    // downloadCurve: function(code) {
    //   window.open("/api/charter/" + code + "/download");
    // },
  },
  computed: {
    productID() {
      return this.id + (this.pnlGroup ? ":" + this.pnlGroup : "");
    },

    pnlClass() {
      if (this.pnl < 0) {
        return "loss";
      }

      if (this.pnl > 0) {
        return "profit";
      }
      return "";
    },
  },
  created: function() {
    this.instanceStatusColor = this.$store.getters.instanceStatusColor;
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
@import "~buefy/src/scss/components/_dropdown.scss";

.header {
  font-size: 12px;
  text-align: left;
  height: 17px;
  line-height: 17px;
}

.product-id {
  display: inline-block;
  color: black;
  margin-left: 2px;
  text-decoration: none;
  width: 45px;
  line-height: 17px;
}

.algo-type {
  display: inline-block;
  color: darkslateblue;
  font-size: 10px;
  padding-left: 4px;
  height: 17px;
  width: 20px;
}

.currency {
  height: 17px;
  color: rgba(0, 0, 0, 0.7);
}

.hovered .currency {
  opacity: 1;
}

.hovered .product-id {
  color: blue;
  text-decoration: underline;
}

.hovered .algo-type {
  color: black;
}

.pnl {
  float: right;
  font-weight: bold;
  padding-right: 2px;
}

.profit {
  color: green;
}

.loss {
  color: red;
}

.hovered .close-button {
  opacity: 1;
}

.close-button {
  font-size: 12px;
  width: 14px;
  height: 17px;
  float: right;
  text-align: center;
  background-color: transparent;
  cursor: pointer;
  color: black;
  opacity: 0.2;
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
