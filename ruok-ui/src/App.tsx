import { lazy, Suspense } from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { Route, Routes } from 'react-router-dom';
import routes from './routes/index';
import Wellcome from './pages/Wellcome';
const DefaultLayout = lazy(() => import('./layouts/default'));

function Loading() {
  return (
    <>
      <div>
        <h1>Loading</h1>
      </div>
    </>
  );
}

function App() {
  const queryClient = new QueryClient();

  return (
    <>
      <QueryClientProvider client={queryClient}>
        <Routes>
          <Route element={<DefaultLayout />}>
            <Route index element={<Wellcome />} />
            {routes.map((r) => {
              return (
                <Route
                  key={r.title}
                  path={r.path}
                  Component={() => (
                    <Suspense fallback={<Loading />}>
                      <r.component />
                    </Suspense>
                  )}
                />
              );
            })}
          </Route>
        </Routes>
      </QueryClientProvider>
    </>
  );
}

export default App;
