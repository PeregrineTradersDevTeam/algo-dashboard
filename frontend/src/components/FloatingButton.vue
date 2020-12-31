<template>
  <div
    class="floating-button"
    @click="$emit('click')"
    :style="{ backgroundColor: bgColor, color: color }"
  >
    {{ title ? title : "" }}
    <span class="icon"
      ><font-awesome-icon v-if="icon" :icon="icon" slot="trigger" />
    </span>
  </div>
</template>

<script>
export default {
  name: "FloatingButton",
  props: {
    title: String,
    icon: String,
    color: {
      type: String,
      default: "rgba(0, 0, 0, 0.7)",
    },
    bgColor: {
      type: String,
    },
  },
  data: function() {
    return {};
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.floating-button {
  right: 1px;
  height: 26px;
  width: 30px;
  cursor: pointer;
  border-radius: 4px;
  margin-left: 0px;
  margin-bottom: 2px;
  margin-top: 2px;
  font-weight: bold;
  font-size: 14px;
  line-height: 25px;
  background-color: silver;
  color: rgba(0, 0, 0, 0.7);
}

.icon {
  display: block;
}

.floating-button:hover {
  background-color: blue;
  color: white;
}
</style>
