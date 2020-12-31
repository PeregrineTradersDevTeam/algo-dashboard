<template>
  <div
    v-bind:class="isActive == true ? 'dark-layer' : ''"
    @keydown.esc.stop="hide"
    tabindex="0"
    ref="modal"
  >
    <div
      v-if="isActive"
      v-bind:class="screen.type"
      v-bind:style="sizes"
      @keyup.esc="hide"
    >
      <div class="window-header" ref="hd" :style="{ backgroundColor: bgColor }">
        <span class="window-title">{{ title }}</span>
        <div v-if="$slots.header" style="float:left; margin:4px 0px 0px 10px">
          <slot name="header"></slot>
        </div>
        <button class="close-button" @click.stop="hide">
          X
        </button>
      </div>
      <div class="window-content">
        <slot></slot>
      </div>
      <div v-if="$slots.footer" class="window-footer">
        <slot name="footer"></slot>
      </div>

      <ResponsiveModalResizeTriangle
        v-if="isActive && !isMobile && resizeable"
        :width="parseInt(rect.width)"
        :height="parseInt(rect.height)"
        @resize="resize"
        @resize-stop="resizeStop"
      />
    </div>
  </div>
</template>

<script>
import ResponsiveModalResizeTriangle from "@/components/ResponsiveModalResizeTriangle.vue";

export default {
  name: "ResponsiveModal",
  components: {
    ResponsiveModalResizeTriangle,
  },
  props: {
    name: {
      type: String,
      default: "",
    },
    isActive: Boolean,
    title: String,
    widthLimit: String,
    resizeable: {
      type: Boolean,
      default: true,
    },
  },
  data: function() {
    return {
      rect: {
        left: "10px",
        top: "10px",
        width: "10px",
        height: "10px",
      },
      isMobile: false,
      screen: {},
      isResizable: {
        type: Boolean,
        default: false,
      },
      moving: false,
      initialX: 0,
      initialY: 0,
      initialLeft: 0,
      initialTop: 0,
      restoreKeyPressed: 0,
    };
  },
  computed: {
    sizes: function() {
      return (
        "left:" +
        this.rect.left +
        ";top:" +
        this.rect.top +
        ";width:" +
        this.rect.width +
        ";height:" +
        this.rect.height +
        ";max-height:" +
        this.rect.height +
        ";"
      );
    },
    border: function() {
      return "width:calc(100% - 2px);height:calc(100% -1px);";
      // (parseInt(this.rect.height, 10) - 36) +
      // "px;"
    },
  },

  methods: {
    resizeStop(event) {
      console.log("resize stop", event);
      if (this.name && !this.isMobile && this.resizeable) {
        localStorage.setItem(this.name, JSON.stringify(this.rect));
      }
    },
    resize(event) {
      let w = event.width;
      let h = event.height;
      if (w < 700) {
        w = 700;
      } else if (w + parseInt(this.rect.left) > window.innerWidth) {
        w = parseInt(this.rect.width);
      }

      if (h < 350) {
        h = 350;
      } else if (h + parseInt(this.top) > window.innerHeight) {
        h = parseInt(this.rect.height);
      }

      this.rect.width = w + "px";
      this.rect.height = h + "px";
    },
    startMove(event) {
      console.log("start move", this.isActive);
      if (!this.isActive) {
        return;
      }

      this.moving = true;
      this.initialX = event.clientX;
      this.initialY = event.clientY;
      this.initialLeft = parseInt(this.rect.left);
      this.initialTop = parseInt(this.rect.top);

      window.addEventListener("mousemove", this.move, false);
      window.addEventListener("mouseup", this.stopMove, false);
      event.stopPropagation();
      event.preventDefault();
    },
    stopMove(event) {
      this.moving = false;
      this.initialX = 0;
      this.initialY = 0;
      window.removeEventListener("mousemove", this.move, false);
      window.removeEventListener("mouseup", this.stopMove, false);
      if (this.name && !this.isMobile) {
        localStorage.setItem(this.name, JSON.stringify(this.rect));
      }
    },
    move(event) {
      event.preventDefault();
      if (!this.moving) {
        return;
      }

      let offset_x = event.clientX - this.initialX;
      let offset_y = event.clientY - this.initialY;

      let top = parseInt(this.initialTop) + offset_y;
      if (top + parseInt(this.rect.height) > window.innerHeight) {
        top = parseInt(this.rect.top);
      }
      if (top < 0) {
        top = 0;
      }

      let left = parseInt(this.initialLeft) + offset_x;
      if (left < 0) {
        left = 0;
      }
      if (left + parseInt(this.rect.width) >= window.innerWidth) {
        left = parseInt(this.rect.left);
      }

      this.rect.top = top + "px";
      this.rect.left = left + "px";
    },
    hide() {
      this.$emit("update:is-active", false);
      if (this.$listeners.close) {
        this.$emit("close");
      }
    },

    restoreSize() {
      this.restoreKeyPressed++;
      if (this.restoreKeyPressed > 3) {
        this.calcSizes();
        this.restoreKeyPressed = 0;
      }
    },

    getScreenProperties() {
      let res = {
        width: window.innerWidth,
        height: window.innerHeight,
        type: "desktop",
        mode: "landscape",
      };

      if (window.orientation == 90 || window.orientation == -90) {
        res.mode = "landscape";
        res.type = "mobile";
      } else {
        if (window.matchMedia("only screen and (max-width: 492px)").matches) {
          res.mode = "portrait";
          res.type = "mobile";
        }
      }
      return res;
    },

    calcSizes() {
      if (this.screen.type == "mobile") {
        this.rect.top = "29px";
        this.rect.left = "0px";
        this.rect.width = "100%";
        this.rect.height = "100%";
        return;
      }
      // desktop

      this.rect.top = 300 / 2 + "px";
      this.rect.height = this.screen.height - 300 + "px";

      if (this.widthLimit.endsWith("%")) {
        let wl = parseInt(this.widthLimit, 10);

        let w = (this.screen.width / 100.0) * wl;
        //console.log("w=", this);
        this.rect.width = w + "px";
        this.rect.left = (this.screen.width - w) / 2 + "px";
        return;
      }

      if (this.widthLimit != "") {
        this.rect.left = this.screen.width / 2 - this.widthLimit / 2 + "px";
        this.rect.width = this.widthLimit + "px";
        return;
      }

      //  widthLimit not assigned
      this.rect.left = 300 / 2 + "px";
      this.rect.width = this.screen.width - 300 + "px";
    },
  },
  watch: {
    isActive: function(n, o) {
      if (this == undefined || n == false) {
        return;
      }
      this.screen = this.getScreenProperties();
      this.isMobile = this.screen.type === "mobile";

      if (this.name === "" || this.screen.type == "mobile") {
        this.calcSizes();
      } else {
        const rect = localStorage.getItem(this.name);
        if (rect) {
          this.rect = JSON.parse(rect);
        } else {
          this.calcSizes();
        }
      }

      if (this.$listeners.show) {
        this.$emit("show", this.screen);
      }

      this.$nextTick(() => {
        if (!this.isMobile) {
          this.$refs.hd.addEventListener(
            "mousedown",
            this.startMove,
            {
              capture: true,
            },
            false
          );
          this.$refs.modal.focus();
        }
      });
    },
  },
  mounted() {
    let that = this;

    // event related to the screen orientation change.
    window.addEventListener("resize", function(event) {
      that.screen = that.getScreenProperties();
      if (that.isActive && that.$listeners.resize) {
        that.$emit("resize", that.screen);
      }
    });
  },

  beforeDestroy: function() {
    window.removeEventListener("resize", this);
  },
};
</script>

<style scoped lang="scss">
.dark-layer {
  position: fixed;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  z-index: 2;
}
.mobile {
  display: block;
  position: fixed;
  background-color: #1e1e1e;
  z-index: 15;
}

.desktop {
  display: block;
  position: absolute;
  background-color: #1e1e1e;
  color: whitesmoke;
}

.window-header {
  display: block;
  background-color: grey;
  font-size: 15pt;
  position: relative;
  width: 100%;
  height: 36px;
}

.window-content {
  position: relative;
  width: calc(100% - 2px);
  height: calc(100% - 36px);
  margin: 0;
  padding: 0;
  border-left: 1px solid gray;
  border-bottom: 1px solid gray;
  border-right: 1px solid gray;
}

.window-title {
  display: inline-block;
  padding: 0px 0px 0px 10px;
  float: left;
  font-size: 22px;
  font-weight: bold;
  color: yellow;
  line-height: 36px;
}

.close-button {
  display: block;
  height: 24px;
  width: 24px;
  line-height: 24px;
  border-radius: 50%;
  border: 1px solid #1e1e1e;
  float: right;
  margin: 5px 10px 0px 10px;
  padding-left: 1px;
}
.close-button:hover {
  background-color: blue;
  color: white;
}

.window-footer {
  position: relative;
  display: block;
  bottom: -1px;
}
</style>
