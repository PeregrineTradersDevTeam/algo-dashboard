<template>
  <div id="app">
    <vue-confirm-dialog></vue-confirm-dialog>
    <!-- <div id="nav">
       <router-link to="/">Home</router-link>|
      <router-link to="/about">About</router-link>
    </div>-->
    <router-view />
  </div>
</template>

<style>
* {
  padding: 0;
  margin: 0;
  font-family: "Dejavue Sans", Tahoma, Arial, sans-serif;
}

a {
  color: blueviolet;
}

a:hover {
  text-decoration: underline;
  cursor: pointer;
}

html {
  margin: 0;
  padding: 0;
}

body {
  font-size: 10px;
  height: 100%;
  background-color: #1e1e1e;
}

#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: silver;
}
/*
#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}*/
</style>
