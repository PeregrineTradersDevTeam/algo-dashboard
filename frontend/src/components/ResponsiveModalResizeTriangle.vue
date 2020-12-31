<template>
  <div class="resize-button">
    <font-awesome-icon icon="expand-alt" flip="vertical" />
  </div>
</template>

<script>
export default {
  name: "ResponsiveModalResizeTriangle",
  props: {
    height: Number,
    width: Number,
  },
  data: function() {
    return {
      clicked: false,
      size: {},
      initialX: 0,
      initialY: 0,
      initialHeight: 0,
      initialWidth: 0,
    };
  },
  computed: {},

  methods: {
    start(event) {
      this.targetClass = event.target.className;
      this.clicked = true;
      this.initialX = event.clientX;
      this.initialY = event.clientY;
      this.initialHeight = this.height;
      this.initialWidth = this.width;
      window.addEventListener("mousemove", this.mousemove, false);
      window.addEventListener("mouseup", this.stop, false);
      event.stopPropagation();
      event.preventDefault();
    },
    stop() {
      this.clicked = false;
      this.targetClass = "";
      this.initialX = 0;
      this.initialY = 0;
      window.removeEventListener("mousemove", this.mousemove, false);
      window.removeEventListener("mouseup", this.stop, false);
      this.$emit("resize-stop");
    },
    mousemove(event) {
      this.emitResize(event);
    },
    emitResize(event) {
      if (!this.clicked) {
        return;
      }
      event.preventDefault();

      let w = this.initialWidth + event.clientX - this.initialX;
      let h = this.initialHeight + event.clientY - this.initialY;

      this.$emit("resize", { width: w, height: h });
    },
  },
  mounted() {
    this.$el.addEventListener("mousedown", this.start, false);
  },
};
</script>

<style scoped>
.resize-button {
  display: block;
  overflow: hidden;
  position: absolute;
  background: transparent;
  z-index: 9999999;
  width: 12px;
  height: 12px;
  bottom: 0;
  right: 0;
  color: rgba(255, 255, 255, 0.5);
}

.resize-button:hover {
  -webkit-box-shadow: inset 0px 0px 4px 0px rgba(9, 13, 237, 1);
  -moz-box-shadow: inset 0px 0px 4px 0px rgba(9, 13, 237, 1);
  box-shadow: 1px 2px 10px 3px slateblue;
}

#vue-modal-triangle::after {
  display: block;
  position: absolute;
  content: "";
  background: transparent;
  left: 0;
  top: 0;
  width: 0;
  height: 0;
  border-bottom: 10px solid #ddd;
  border-left: 10px solid transparent;
}

#vue-modal-triangle.clicked::after {
  border-bottom: 10px solid #369be9;
}
</style>
