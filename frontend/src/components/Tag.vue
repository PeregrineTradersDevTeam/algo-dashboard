<template>
  <div class="tag">
    <div
      :class="[short ? 'label hideable' : 'label']"
      :style="{ 'background-color': labelBC }"
    >
      <span :style="{ color: labelColor }">{{ label }}</span>
    </div>
    <div
      class="value"
      :style="{ 'background-color': valueBC, color: valueColor }"
    >
      {{ value }}
    </div>
  </div>
</template>

<script>
export default {
  name: "Tag",
  props: {
    short: {
      default: false,
      type: Boolean,
    },
    label: String,
    labelColor: {
      type: String,
      default: "silver",
    },
    labelBC: {
      type: String,
      default: "#4a4a4a",
    },
    value: [String, Number],
    valueBC: {
      type: String,
      default: "#ffdd57",
    },
    valueColor: {
      type: String,
      default: "#fff",
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.tag {
  display: inline-flex;
  margin: 0;
  margin-right: 5px;
  font-size: 12px;
  height: 2px;
}

.hideable {
  display: none;
}

.label {
  padding: 2px 5px 2px 5px;
  height: 20px;
  vertical-align: middle;
  line-height: 20px;
  border-radius: 4px 0px 0px 4px;
  flex-shrink: 0;
}
@media (min-width: 900px) {
  .hideable {
    display: inline-block;
  }
}

.value {
  display: inline-block;
  padding: 2px 5px 2px 5px;
  height: 20px;
  vertical-align: middle;
  line-height: 20px;
  border-radius: 0px 4px 4px 0px;
}
</style>
