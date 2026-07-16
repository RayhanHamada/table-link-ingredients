import { Routes, Route, Navigate, NavLink } from "react-router-dom";
import { Toaster } from "sonner";
import { cn } from "@/lib/utils";

import { IngredientList } from "@/features/ingredients/ingredient-list";
import { IngredientForm } from "@/features/ingredients/ingredient-form";
import { ItemList } from "@/features/items/item-list";
import { ItemForm } from "@/features/items/item-form";

function Layout({ children }: { children: React.ReactNode }) {
  const linkClasses = ({ isActive }: { isActive: boolean }) =>
    cn(
      "inline-flex items-center rounded-lg px-3 py-2 text-sm font-medium transition-colors",
      isActive
        ? "bg-muted text-foreground"
        : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
    );

  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-40 border-b bg-background">
        <div className="flex h-14 items-center gap-6 px-6">
          <NavLink to="/ingredients" className="font-heading text-lg font-semibold tracking-tight">
            TableLink
          </NavLink>
          <nav className="flex items-center gap-1">
            <NavLink to="/ingredients" className={linkClasses}>
              Ingredients
            </NavLink>
            <NavLink to="/items" className={linkClasses}>
              Items
            </NavLink>
          </nav>
        </div>
      </header>
      <main className="flex-1 px-6 py-8">
        <div className="mx-auto max-w-5xl">{children}</div>
      </main>
    </div>
  );
}

export default function App() {
  return (
    <>
      <Routes>
        <Route
          path="/"
          element={
            <Layout>
              <Navigate to="/ingredients" replace />
            </Layout>
          }
        />
        <Route
          path="/ingredients"
          element={
            <Layout>
              <IngredientList />
            </Layout>
          }
        />
        <Route
          path="/ingredients/new"
          element={
            <Layout>
              <IngredientForm />
            </Layout>
          }
        />
        <Route
          path="/ingredients/:uuid/edit"
          element={
            <Layout>
              <IngredientForm />
            </Layout>
          }
        />
        <Route
          path="/items"
          element={
            <Layout>
              <ItemList />
            </Layout>
          }
        />
        <Route
          path="/items/new"
          element={
            <Layout>
              <ItemForm />
            </Layout>
          }
        />
        <Route
          path="/items/:uuid/edit"
          element={
            <Layout>
              <ItemForm />
            </Layout>
          }
        />
        <Route
          path="*"
          element={
            <Layout>
              <Navigate to="/ingredients" replace />
            </Layout>
          }
        />
      </Routes>
      <Toaster />
    </>
  );
}

