<template>
  <ul class="tabs">
    <li v-for="(item, index) in title" :key="index">
      <input
        type="radio"
        name="tabs"
        :id="'tab' + index"
        :checked="index == 0"
      />
      <label :for="'tab' + index" role="tab" :tabindex="index">{{
        item
      }}</label>
      <div :id="'tab-content' + index" class="tab-content" role="tabpanel">
        <slot v-bind:name="item"></slot>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  name: "Tabs",
  props: {
    title: Array,
    components: Array,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
//@import "compass/css3";

//@import url("https://fonts.googleapis.com/css?family=Lato");

//$background: #9b59b6;
$background: #9b59b6;
$tabs-base-color: #8e44ad;
// * {
//   margin: 0;
//   padding: 0;
//   // @include box-sizing(border-box);
// }
// body {
//   text-align: left;
//   color: #fff;
//   background: $background;
// }

.tabs {
  display: flex;
  flex-direction: row;
  width: calc(100% - 10px);
  list-style: none;
  text-align: left;
  padding: 5px;
  li {
    float: left;
    display: block;
  }
  input[type="radio"] {
    position: absolute;
    top: 0;
    left: -9999px;
  }
  label {
    display: block;
    padding: 10px 21px;
    border-radius: 2px 2px 0 0;
    font-size: 14px;
    font-weight: normal;
    text-transform: uppercase;
    background: $tabs-base-color;
    cursor: pointer;
    position: relative;
    top: 4px;
    //   @include transition(all 0.2s ease-in-out);
    &:hover {
      background: darken($tabs-base-color, 10);
    }
  }
  .tab-content {
    z-index: 2;
    display: none;
    overflow: hidden;
    width: 100%;
    font-size: 17px;
    line-height: 25px;
    position: absolute;
    top: 42px;
    left: 0;
    background: transparent; //darken($tabs-base-color, 15);
    border-top: 7px solid darken($tabs-base-color, 15);
    height: calc(100% - 63px);
  }

  //The Magic
  [id^="tab"]:checked + label {
    top: 0;
    padding-top: 17px;
    background: darken($tabs-base-color, 15);
  }
  [id^="tab"]:checked ~ [id^="tab-content"] {
    display: flex;
  }
}
p.link {
  clear: both;
  margin: 380px 0 0 15px;
  a {
    text-transform: uppercase;
    text-decoration: none;
    display: inline-block;
    color: #fff;
    padding: 5px 10px;
    margin: 0 5px;
    background-color: darken($tabs-base-color, 15);
    //  @include transition(all 0.2s ease-in);
    &:hover {
      background-color: darken($tabs-base-color, 20);
    }
  }
}
</style>
