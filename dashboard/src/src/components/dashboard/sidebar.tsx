"use client"

import { useState } from "react"
import { NavLink } from "react-router-dom"
import { Users, Tag, Settings, LogOut, Menu, X, Clock } from "lucide-react"
import { cn } from "../../lib/utils"
import { Button } from "../ui/button"
import { Sheet, SheetContent } from "../ui/sheet"
import { useIsMobile } from "../../hooks/use-mobile"
import { logout } from "@/services/auth-service"

interface SidebarProps {
  className?: string
}

export function Sidebar({ className }: SidebarProps) {
  const [isOpen, setIsOpen] = useState(true)
  const [isMobileOpen, setIsMobileOpen] = useState(false)
  const isMobile = useIsMobile()

  const toggleSidebar = () => {
    if (isMobile) {
      setIsMobileOpen(!isMobileOpen)
    } else {
      setIsOpen(!isOpen)
    }
  }

  const sidebarItems = [
    {
      title: "ユーザー管理",
      href: "/users",
      icon: Users,
    },
    {
      title: "ラベル管理",
      href: "/labels",
      icon: Tag,
    },
    {
      title: "プロバイダ設定",
      href: "/providers",
      icon: Settings,
    },
    {
      title: "セッション管理",
      href: "/sessions",
      icon: Clock,
    },
  ]

  const SidebarContent = () => (
    <div className="flex h-full flex-col">
      <div className="flex h-14 items-center border-b px-4">
        <h2 className="text-lg font-semibold">認証基盤ダッシュボード</h2>
        {isMobile && (
          <Button variant="ghost" size="icon" className="ml-auto" onClick={() => setIsMobileOpen(false)}>
            <X className="h-5 w-5" />
          </Button>
        )}
      </div>
      <div className="flex-1 overflow-auto py-2">
        <nav className="grid gap-1 px-2">
          {sidebarItems.map((item) => (
            <NavLink
              key={item.href}
              to={item.href}
              className={({ isActive }) =>
                cn(
                  "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground",
                  isActive ? "bg-accent text-accent-foreground" : "transparent",
                )
              }
            >
              <item.icon className="h-5 w-5" />
              <span className={cn(isOpen ? "opacity-100" : "opacity-0", "transition-opacity duration-200")}>
                {item.title}
              </span>
            </NavLink>
          ))}
        </nav>
      </div>
      <div className="border-t p-4">
        <Button variant="outline" className="w-full justify-start gap-2" onClick={() => {
          logout();
          window.location.reload();
        }}>
          <LogOut className="h-4 w-4" />
          <span className={cn(isOpen ? "opacity-100" : "opacity-0", "transition-opacity duration-200")}>
            ログアウト
          </span>
        </Button>
      </div>
    </div>
  )

  // モバイル用のサイドバー
  if (isMobile) {
    return (
      <>
        <Button variant="ghost" size="icon" className="fixed left-4 top-4 z-40 md:hidden" onClick={toggleSidebar}>
          <Menu className="h-5 w-5" />
        </Button>
        <Sheet open={isMobileOpen} onOpenChange={setIsMobileOpen}>
          <SheetContent side="left" className="p-0">
            <SidebarContent />
          </SheetContent>
        </Sheet>
      </>
    )
  }

  // デスクトップ用のサイドバー
  return (
    <div
      className={cn(
        "relative flex h-screen flex-col border-r bg-background transition-all duration-300 ease-in-out",
        isOpen ? "w-64" : "w-[70px]",
        className,
      )}
    >
      {/* <Button
        variant="ghost"
        size="icon"
        className="absolute -right-3 top-6 z-10 h-6 w-6 rounded-full border bg-background"
        onClick={toggleSidebar}
      >
        {isOpen ? <span className="text-xs">◀</span> : <span className="text-xs">▶</span>}
      </Button> */}
      <SidebarContent />
    </div>
  )
}
