<template>
  <div class="tab-content-console">
    <div class="subheader">
      <ResponsiveButton title="Refresh" :click="getOrders(code)" />
    </div>
    <div class="aa">
      <Table
        :rows="items"
        :expanding="true"
        :head="['Date', 'Time', 'Order details']"
        :colstyle="{
          0: 'width: 75px',
          1: 'width: 80px',
        }"
      />
    </div>
  </div>
</template>

<script>
async function fetchOrders(id) {
  let url = "/api/log/" + id + "/orders";

  const resp = await fetch(url);
  var data = await resp.json();

  let res = [];
  for (var i in data) {
    res.push({
      atd: formatD(data[i].at),
      att: formatT(data[i].at),
      msg: data[i].msg,
    });
  }

  return res;
}
import Table from "@/components/Table.vue";
import { formatT, formatD } from "@/utils/formatters.js";
import ResponsiveButton from "@/components/ResponsiveButton.vue";

export default {
  name: "InstanceModalOrders",
  components: {
    Table,
    ResponsiveButton,
  },
  props: {
    code: String,
  },
  data: function() {
    return {
      items: [],
    };
  },
  methods: {
    getOrders(id) {
      let that = this;
      fetchOrders(id).then((arr) => {
        that.items = [];
        that.items.push(...arr);
      });
    },
  },
  mounted() {
    this.getOrders(this.code);
  },
};
</script>

<style scoped>
.subheader {
  display: flex;
  padding: 4px;
  position: sticky;
  top: 0px;
  background: #1e1e1e;
  height: 25px;
  justify-content: space-between;
  font-size: 120%;
}

.tab-content-console {
  display: block;
  width: 100%;
  height: calc(100% - 2px);
  background-color: #1e1e1e;
  overflow-y: hidden;
}

.footer {
  display: flex;
  position: sticky;
  bottom: 0px;
  background: #1e1e1e;
  height: 25px;
  justify-content: space-between;
  font-size: 120%;
}
.aa {
  display: block;
  height: calc(100% - 40px);
  top: 0px;
  margin: 0;
  background-color: #1e1e1e;
}
</style>
