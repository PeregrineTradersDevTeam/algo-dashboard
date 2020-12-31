<template>
  <responsive-modal
    name="window-market-holidays"
    :isActive.sync="isActive"
    title="Market Holidays"
    widthLimit="800"
  >
    <div class="holiday-content">
      <div class="subheader">
        <div>
          <p class="last-holiday-label">
            Server date:
            <span class="last-holiday-value">{{ holidays.today }}</span>
          </p>
        </div>
        <div>
          <p class="last-holiday-label">
            <span style="float:right">
              Last day:
              <span class="last-holiday-value">{{ holidays.lastDay }}</span>
            </span>
          </p>
        </div>
      </div>
      <Table
        :rows="holidays.dates"
        :head="['Date', 'Name', 'Market', 'Open time', 'Close', 'MIC']"
        :selection="{
          func: function(item) {
            return item.dt == holidays.today;
          },
          style: 'background-color: red; color:white',
        }"
        :colstyle="{ 0: 'min-width: 73px' }"
      />
    </div>
  </responsive-modal>
</template>

<script>
import ResponsiveModal from "@/components/ResponsiveModal.vue";
import Table from "@/components/Table.vue";

export default {
  name: "MarketHolidaysModal",
  components: {
    "responsive-modal": ResponsiveModal,
    Table: Table,
  },
  data: function() {
    return {
      dates: [],
      today: "",
      isActive: false,
      sortField: "dt",
      sortOrder: "asc",
      serverTime: "",
      lastDay: "",
      isHolidayToday: false,
      isUpdateRequired: false,
      isUserNotified: false,
    };
  },
  computed: {
    // todays() {
    //   let res = [];
    //   let h = this.holidays();
    //   for (let i in h.dates) {
    //     if (h.dates[i].dt == h.today) {
    //       res.push(h.dates[i]);
    //     }
    //   }
    //   return res;
    // },
    holidays() {
      return this.$store.getters.holidays;
    },
  },
  methods: {
    iconColor() {
      return this.isUpdateRequired ? "orange" : "silver";
    },
    backColor() {
      this.isHolidayToday != ""
        ? "background-color: red"
        : "background-color: gray";
    },
    showModal() {
      this.isActive = true;
    },
    // getHolidays: function() {
    //   let that = this;
    //   fetch("/api/holidays")
    //     .then((r) => r.json())
    //     .then(function(j) {
    //       that.today = j.today;
    //       that.lastDay = "";
    //       that.isUpdateRequired = j.isUpdateRequired;
    //       that.isHolidayToday = j.isHolidayToday;
    //       that.$store.commit("setHoliday", j.isHolidayToday);
    //       that.serverTime = j.serverTime;

    //       if (j.dates == null) {
    //         return;
    //       }

    //       if (j.dates.length == 0) {
    //         return;
    //       }
    //       that.dates.push(...j.dates);
    //       that.lastDay = j.dates[j.dates.length - 1].dt;
    //       if (
    //         that.isHolidayToday &&
    //         that.isUserNotified == false &&
    //         !window.webpackHotUpdate // do not show in the dev mode
    //       ) {
    //         that.isUserNotified = true;
    //         that.isActive = true;
    //       }
    //     });
    // },
  },

  mounted: function() {
    // let that = this;
    // this.pooling = setInterval(function() {
    //   that.getHolidays();
    // }, 60000);
  },
  beforeDestroy() {
    // clearInterval(this.pooling);
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
.button-calendar {
  height: 24px;
  width: 28px;
  padding: 1px 1px 0px 1px;
  font-weight: bold;
  font-size: 1.2em;
  border: 1px solid silver;
  cursor: pointer;
  border-radius: 4px;
  margin-left: 0px;
}

.button-calendar:hover {
  background-color: blue;
}

.last-holiday-label {
  margin: 5px;
  text-align: left;
  color: white;
  font-size: 1em;
}

.last-holiday-value {
  margin: 5px;
  text-align: left;
  color: yellow;
  font-weight: bold;
}
/* 
.table {
  background-color: #3a3f44;
  color: #aaa;
}

th {
  vertical-align: top;
  color: #fff;
}

.table td,
.table th {
  border: 1px solid #52575c;
  border-width: 0 0 1px;
  padding: 0.5em 0.75em;
  vertical-align: top;
  color: #fff;
}
.table td.is-primary,
.table th.is-primary {
  background-color: #52575c;
  border-color: #52575c;
  color: #fff;
}

.table td.is-narrow,
.table th.is-narrow {
  white-space: nowrap;
  width: 1%;
}

.table th {
  color: white;
}
.table th:not([align]) {
  text-align: left;
}
.table thead {
  background-color: transparent;
}
.table thead td,
.table thead th {
  border-width: 0 0 1px;
  color: #9d9d9d;
}
.table tfoot {
  background-color: transparent;
}
.table tfoot td,
.table tfoot th {
  border-width: 2px 0 0;
  color: #9d9d9d;
}
.table tbody {
  background-color: transparent;
}
.table tbody tr:last-child td,
.table tbody tr:last-child th {
  border-bottom-width: 0;
}

.table.is-hoverable tbody tr:not(.is-selected):hover {
  background-color: #272b30;
}
.table.is-hoverable.is-striped tbody tr:not(.is-selected):hover {
  background-color: #272b30;
}
.table.is-hoverable.is-striped
  tbody
  tr:not(.is-selected):hover:nth-child(even) {
  background-color: #30353b;
}
.table.is-narrow td,
.table.is-narrow th {
  padding: 0.25em 0.5em;
}

.table.is-striped tbody tr:not(.is-selected):nth-child(even) {
  background-color: #272b30;
} */

/* .table-container {
  -webkit-overflow-scrolling: touch;
  overflow: auto;
  overflow-y: hidden;
  max-width: 100%;
} */

.subheader {
  display: flex;
  position: sticky;
  top: 0px;
  background: #1e1e1e;
  height: 25px;
  justify-content: space-between;
  font-size: 120%;
}

.icon {
  display: block;
  margin-top: 3px;
}

.table-body {
}

.holiday-content {
  display: block;
  height: calc(100% - 35px);
  top: 0px;
  margin: 0;
  background-color: #1e1e1e;
}
</style>
