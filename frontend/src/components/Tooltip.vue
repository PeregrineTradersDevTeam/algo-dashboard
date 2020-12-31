<template>
  <div class="tooltip-box">
    <div class="tooltip">
      <span class="triangle"></span>
    </div>
  </div>
</template>

<script>
export default {
  name: "Tooltip",
  props: {
    msg: String,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.tooltip {
  width: 150px;
  background: #59c7f9;
  color: #ffffff;
  text-align: center;
  padding: 10px 20px 10px 20px;
  border-radius: 5px;
  bottom: calc(100% + 14px);
  left: 50%;
  transform: translate-x(-50%);
  z-index: 6;
}
.tooltip-box {
  position: relative;
}
.tooltip-box :hover {
}
.triangle {
  border-width: 0 7.5px 9px;
  border-color: transparent;
  border-bottom-color: #59c7f9;
  position: absolute;
  bottom: -9px;
  transform: rotate(180deg) translate-x(-50%);
  left: calc(50% - 15px);
}
</style>
