<template>
  <div class="floating-button" @click="$emit('click')">
    {{ title ? title : "" }}
    <font-awesome-icon v-if="icon" :icon="icon" slot="trigger" />
  </div>
</template>

<script>
export default {
  name: "SortButton",
  props: {
    title: String,
    icon: String,
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
  height: 24px;
  width: 28px;
  border: 1px solid silver;
  cursor: pointer;
  border-radius: 4px;
  margin-left: 0px;
  background-color: silver;
  margin-bottom: 2px;
  margin-top: 2px;
  font-weight: bold;
  font-size: 10px;
  vertical-align: middle;
  padding-top: 4px;
}

.floating-button:hover {
  background-color: blue;
  color: white;
}
</style>
