<template>
  <div>
    <span @click="toggleShow" class="anchor">{{ title }}</span>
    <div v-if="showMenu" class="menu" id="instance-dropdown-menu">
      <div
        class="menu-item"
        v-for="(item, index) in items"
        :key="index"
        @mouseup="itemClicked"
      >
        {{ item }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    title: String,
    items: {
      type: Array,
    },
  },
  data: function() {
    return {
      showMenu: false,
    };
  },
  methods: {
    toggleShow: function() {
      if (!this.showMenu) {
        window.addEventListener("mouseup", this.outsideClickCapturer);
      } else {
        window.removeEventListener("mouseup", this.outsideClickCapturer);
      }
      this.showMenu = !this.showMenu;
    },

    outsideClickCapturer(e) {
      console.log("click", e);
      this.toggleShow();
      // if (
      //   !document.getElementById("l2").contains(e.target) &&
      //   !document.getElementById("logo-menu").contains(e.target)
      // ) {
      //   alert("Clicked outside l2 and logo-menu");
      //   document.getElementById("l2").style.height = "0px"; //the same code you've used to hide the menu
      // }
    },
    itemClicked: function() {
      this.toggleShow();
      this.$emit("click");
    },
  },
};
</script>

<style scoped>
.anchor {
  display: inline;
  /* align-items: center;
  /* justify-content: center; */
  border: 1px solid transparent;
  /* padding: 0.75rem 2rem; */
  /* font-size: 1rem; */
  /* border-radius: 0.25rem; */
  /* transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out,
    border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out; */
  color: #000;
  /* background-color: #27ae60; */
  background-color: transparent;
  /* border-color: #27ae60; */
}

.anchor::after {
  display: inline-block;
  width: 0;
  height: 0;
  margin-left: 0;
  vertical-align: 0.255em;
  content: "";
  border-top: 0.3em solid;
  border-right: 0.28em solid transparent;
  border-bottom: 0;
  border-left: 0.28em solid transparent;
}

.anchor:hover {
  color: #fff;
  background-color: blue;
  /* border-color: #229954; */
  cursor: pointer;
}

.menu {
  background-color: #fff;
  background-clip: padding-box;
  /* border: 1px solid rgba(0, 0, 0, 0.15); */
  border-radius: 0.25rem;
  color: #212529;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  /* font-size: 1rem; */
  list-style: none;
  margin: 0.125rem 0 0;
  padding: 0.5rem 0;
  position: absolute;
  text-align: left;
  z-index: 13;
}

.menu-item {
  color: #212529;
  padding: 0.25rem 1.5rem;
  transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out,
    border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.menu-item:hover {
  background-color: blue;
  color: white;
  cursor: pointer;
}
</style>
