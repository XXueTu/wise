import { Dashboard } from "@/components/Dashboard"
import { Layout } from "@/components/Layout"
import { ModelManager } from "@/components/ModelManager"
import { ResourceManager } from "@/components/ResourceManager"
import { TagManager } from "@/components/TagManager"
import { TaskManager } from "@/components/TaskManager"
import { Route, BrowserRouter as Router, Routes } from "react-router-dom"

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Dashboard />} />
          <Route path="resources" element={<ResourceManager />} />
          <Route path="models" element={<ModelManager />} />
          <Route path="tasks" element={<TaskManager />} />
          <Route path="tags" element={<TagManager />} />
        </Route>
      </Routes>
    </Router>
  )
}

export default App
