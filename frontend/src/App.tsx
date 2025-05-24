import { Dashboard } from "@/components/Dashboard"
import { ModelManager } from "@/components/ModelManager"
import { ResourceManager } from "@/components/ResourceManager"
import { Button } from "@/components/ui/button"
import { useState } from "react"

function App() {
  const [activeTab, setActiveTab] = useState<"dashboard" | "resources" | "models">("dashboard")

  return (
    <div className="min-h-screen bg-background">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 items-center">
          <div className="mr-4 flex">
            <a
              className="mr-6 flex items-center space-x-2 cursor-pointer"
              onClick={() => setActiveTab("dashboard")}
            >
              <span className="font-bold">管理系统</span>
            </a>
            <nav className="flex items-center space-x-6 text-sm font-medium">
              <Button
                variant={activeTab === "resources" ? "default" : "ghost"}
                onClick={() => setActiveTab("resources")}
              >
                资源管理
              </Button>
              <Button
                variant={activeTab === "models" ? "default" : "ghost"}
                onClick={() => setActiveTab("models")}
              >
                模型管理
              </Button>
            </nav>
          </div>
        </div>
      </header>
      <main>
        {activeTab === "dashboard" && <Dashboard />}
        {activeTab === "resources" && <ResourceManager />}
        {activeTab === "models" && <ModelManager />}
      </main>
    </div>
  )
}

export default App
