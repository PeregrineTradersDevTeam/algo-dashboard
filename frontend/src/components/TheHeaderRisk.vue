<template>
  <section>
    <Tag
      label="pDelta"
      :value="lastPd"
      labelBC="rgb(114 77 129)"
      valueBC="lightcoral"
    />
    <Tag
      label="pDLimit"
      :value="lastPl"
      labelBC="rgb(114 77 129)"
      valueBC="lightcoral"
    />
    <Tag
      label="pDLCash"
      :value="pdlc"
      labelBC="rgb(114 77 129)"
      valueBC="lightcoral"
    />
  </section>
</template>

<script>
import Tag from "@/components/Tag.vue";
export default {
  name: "TheHeaderRisk",
  components: {
    Tag,
  },
  data: function() {
    return {};
  },
  computed: {
    lastPd() {
      return this.$store.state.lastPd;
    },
    lastPl() {
      return this.$store.state.lastPl;
    },
    pdlc() {
      return this.$store.state.pdlc;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.page-header {
  /* z-index: 10px; */
}
</style>
