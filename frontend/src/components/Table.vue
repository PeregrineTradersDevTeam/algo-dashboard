<template>
  <div class="table-container">
    <table width="100%" border="0">
      <thead v-if="head">
        <tr>
          <td v-if="expanding" class="th"></td>
          <th
            class="th"
            v-for="(val, index) in head"
            :key="index"
            :style="colstyle && colstyle[index] ? colstyle[index] : ''"
          >
            {{ val }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          class="tr"
          v-for="(row, index) in rows"
          :key="index"
          :style="selection && selection.func(row) ? selection.style : ''"
        >
          <td v-if="expanding" class="table-cell">[+]</td>
          <td
            class="table-cell"
            v-for="(val, attr, index) in row"
            :key="attr"
            :style="tdstyle(index, val)"
          >
            {{ replacer && replacer[index] ? replacer[index](val) : val }}
          </td>
        </tr>
      </tbody>
    </table>
    <slot name="footer"></slot>
  </div>
</template>
//
<!-- :style="colstyle && colstyle[index] ? colstyle[index] : ''" -->

<script>
export default {
  name: "Table",
  props: {
    expanding: Boolean,
    head: Array,
    rows: Array,
    selection: Object,
    colstyle: Object,
    cellstyle: Object,
    replacer: Object,
  },
  methods: {
    tdstyle(index, val) {
      if (!this.cellstyle) {
        return "";
      }
      if (!this.cellstyle[index]) {
        return "";
      }
      if (typeof this.cellstyle[index] === "string") {
        return this.cellstyle[index];
      }
      if (typeof this.cellstyle[index] === "function") {
        return this.cellstyle[index](val);
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
$color1: #30353b;
$color2: transparent;

.table-container {
  position: relative;
  width: 100%;
  text-align: left;
  overflow: auto;
  height: calc(100% - 1px);
  font-size: 13px;
  color: silver;
}

table {
  text-align: left;
  position: relative;
  width: 100%;
  border-collapse: collapse; // allows using "border:0" on <tr> level
  overflow: auto;
  max-height: calc(100% - 1px);
}

.table-cell {
  padding: 2px;
  vertical-align: top;
}

.tr {
  background-color: $color1;
  border: 0;
}
.tr:nth-child(even) {
  background-color: $color2;
}

.th {
  background: rgb(16, 14, 84);
  position: sticky;
  top: 0px;
  padding: 4px;
}
</style>
