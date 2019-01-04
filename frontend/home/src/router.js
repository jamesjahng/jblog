import Vue from 'vue';
import Router from 'vue-router';

import About from './components/About.vue'
import Blog from './components/Blog.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'About',
      component: About
    },
    {
      path: '/blog',
      name: 'Blog',
      component: Blog,
      children: [
        {
          path: ':year/:month/:day/:title'
        }
      ]
    } 
  ]
});
