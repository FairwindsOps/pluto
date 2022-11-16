// To see all options:
// https://vuepress.vuejs.org/config/
// https://vuepress.vuejs.org/theme/default-theme-config.html
module.exports = {
  title: "Pluto Documentation",
  description: "Documentation for Fairwinds' Pluto",
  themeConfig: {
    docsRepo: "FairwindsOps/pluto",
    sidebar: [
      {
        title: "Pluto",
        path: "/",
        sidebarDepth: 0,
      },
      {
        title: "Installation",
        path: "/installation",
      },
      {
        title: "Quickstart",
        path: "/quickstart",
      },
      {
        title: "Advanced Usage",
        path: "/advanced",
      },
      {
        title: "FAQ",
        path: "/faq",
      },
      {
        title: "CircleCI Orb",
        path: "/orb",
      },
      {
        title: "Contributing",
        children: [
          {
            title: "Guide",
            path: "contributing/guide"
          },
          {
            title: "Code of Conduct",
            path: "contributing/code-of-conduct"
          }
        ]
      }
    ]
  },
}
