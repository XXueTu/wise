import { Button } from "@/components/ui/button"
import { Outlet, useLocation, useNavigate } from "react-router-dom"

export function Layout() {
  const navigate = useNavigate()
  const location = useLocation()

  return (
    <div className="min-h-screen bg-background">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 items-center">
          <div className="mr-4 flex">
            <a
              className="mr-6 flex items-center space-x-2 cursor-pointer"
              onClick={() => navigate("/")}
            >
              <span className="font-bold">WISE</span>
            </a>
            <nav className="flex items-center space-x-6 text-sm font-medium">
              <Button
                variant={location.pathname === "/" ? "default" : "ghost"}
                onClick={() => navigate("/")}
              >
                仪表盘
              </Button>
              <Button
                variant={location.pathname === "/resources" ? "default" : "ghost"}
                onClick={() => navigate("/resources")}
              >
                资源管理
              </Button>
              <Button
                variant={location.pathname === "/models" ? "default" : "ghost"}
                onClick={() => navigate("/models")}
              >
                模型管理
              </Button>
              <Button
                variant={location.pathname === "/tags" ? "default" : "ghost"}
                onClick={() => navigate("/tags")}
              >
                标签管理
              </Button>
            </nav>
          </div>
        </div>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  )
} 