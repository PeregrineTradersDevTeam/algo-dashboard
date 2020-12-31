<template>
  <div class="tab-content-info">
    <div class="aa">
      <Table :rows="info" :head="['Param', 'Value']"> </Table>
    </div>
  </div>
</template>

<script>
import Table from "@/components/Table.vue";
export default {
  name: "InstanceModalInfo",
  components: {
    Table,
  },
  props: {
    code: String,
  },
  data: function() {
    return {
      info: [],
    };
  },

  mounted() {
    let that = this;
    fetch(`/api/instance/${this.code}/info`)
      .then((r) => r.json())
      .then(function(j) {
        for (let key in j) that.info.push({ k: key, v: j[key] });
      });
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.tab-content-info {
  display: block;
  width: 100%;
  height: calc(100% - 2px);
  background-color: #1e1e1e;
  overflow-y: hidden;
}
</style>
