<template>
  <div
    class="responsive-button"
    :style="{ backgroundColor: bgColor, color: color }"
    @click="
      {
        if (disabled) {
          return;
        }
        $emit('click', $event);
      }
    "
  >
    <p style="text-align: center">
      <strong>{{ title[0] }}</strong
      >{{ !isShort ? title.substring(1) : "" }}
    </p>
  </div>
</template>

<script>
export default {
  name: "ResponsiveButton",
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
    title: String,
    isShort: Boolean,
    color: String,
    bgColor: String,
    minw: {
      type: Number,
    },
  },
};
</script>

<style scoped>
.responsive-button {
  display: block;
  height: 24px;
  padding: 0 4px 0 4px;
  cursor: pointer;
  border-radius: 3px;
  background-color: silver;
  font-size: 12px;
  line-height: 24px;
  margin: inherit;
  color: rgba(0, 0, 0, 0.7);
  white-space: nowrap;
  min-width: 18px;
}

.responsive-button:hover {
  background-color: blue;
  color: white;
}

.label {
  display: flex;
  justify-content: center;
}

.hideable {
  display: none;
}
</style>
