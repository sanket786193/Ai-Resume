import { lazy, Suspense } from 'react'
import { createBrowserRouter, Navigate, RouterProvider, Outlet } from 'react-router-dom'
import { ROUTES, ROLES } from '@/constants'
import MainLayout from '@/components/layout/MainLayout'
import AuthLayout from '@/components/layout/AuthLayout'
import { ProtectedRoute } from '@/components/layout/ProtectedRoute'

const LoginPage = lazy(() => import('@/app/auth/LoginPage').then((m) => ({ default: m.LoginPage })))
const RegisterPage = lazy(() => import('@/app/auth/RegisterPage').then((m) => ({ default: m.RegisterPage })))
const HrDashboardPage = lazy(() => import('@/app/hr/HrDashboardPage').then((m) => ({ default: m.HrDashboardPage })))
const HrJobsPage = lazy(() => import('@/app/hr/HrJobsPage').then((m) => ({ default: m.HrJobsPage })))
const HrJobApplicantsPage = lazy(() => import('@/app/hr/HrJobApplicantsPage').then((m) => ({ default: m.HrJobApplicantsPage })))
const HrAtsPage = lazy(() => import('@/app/hr/HrAtsPage').then((m) => ({ default: m.HrAtsPage })))
const HrInterviewsPage = lazy(() => import('@/app/hr/HrInterviewsPage').then((m) => ({ default: m.HrInterviewsPage })))
const HrOffersPage = lazy(() => import('@/app/hr/HrOffersPage').then((m) => ({ default: m.HrOffersPage })))
const CandidateDashboardPage = lazy(() => import('@/app/candidate/CandidateDashboardPage').then((m) => ({ default: m.CandidateDashboardPage })))
const CandidateApplicationsPage = lazy(() => import('@/app/candidate/CandidateApplicationsPage').then((m) => ({ default: m.CandidateApplicationsPage })))
const PublicJobsPage = lazy(() => import('@/app/public/PublicJobsPage').then((m) => ({ default: m.PublicJobsPage })))
const PublicJobDetailPage = lazy(() => import('@/app/public/PublicJobDetailPage').then((m) => ({ default: m.PublicJobDetailPage })))
const PublicApplyPage = lazy(() => import('@/app/public/PublicApplyPage').then((m) => ({ default: m.PublicApplyPage })))

function LazyRoute({ Component }) {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center text-muted-foreground">Loading...</div>}>
      <Component />
    </Suspense>
  )
}

const router = createBrowserRouter([
  {
    path: ROUTES.HOME,
    element: <MainLayout />,
    children: [
      { index: true, element: <Navigate to={ROUTES.PUBLIC_JOBS} replace /> },
      {
        path: 'login',
        element: <AuthLayout />,
        children: [{ index: true, element: <LazyRoute Component={LoginPage} /> }],
      },
      {
        path: 'register',
        element: <AuthLayout />,
        children: [{ index: true, element: <LazyRoute Component={RegisterPage} /> }],
      },
      {
        path: 'hr',
        element: (
          <ProtectedRoute allowedRoles={[ROLES.HR]}>
            <Outlet />
          </ProtectedRoute>
        ),
        children: [
          { index: true, element: <LazyRoute Component={HrDashboardPage} /> },
          { path: 'jobs', element: <LazyRoute Component={HrJobsPage} /> },
          { path: 'jobs/new', element: <LazyRoute Component={HrJobsPage} /> },
          { path: 'jobs/:jobId/edit', element: <LazyRoute Component={HrJobsPage} /> },
          { path: 'jobs/:jobId/applicants', element: <LazyRoute Component={HrJobApplicantsPage} /> },
          { path: 'ats', element: <LazyRoute Component={HrAtsPage} /> },
          { path: 'interviews', element: <LazyRoute Component={HrInterviewsPage} /> },
          { path: 'offers', element: <LazyRoute Component={HrOffersPage} /> },
        ],
      },
      {
        path: 'candidate',
        element: (
          <ProtectedRoute allowedRoles={[ROLES.CANDIDATE]}>
            <Outlet />
          </ProtectedRoute>
        ),
        children: [
          { index: true, element: <LazyRoute Component={CandidateDashboardPage} /> },
          { path: 'applications', element: <LazyRoute Component={CandidateApplicationsPage} /> },
        ],
      },
      {
        path: 'jobs',
        element: <Outlet />,
        children: [
          { index: true, element: <LazyRoute Component={PublicJobsPage} /> },
          { path: ':jobId', element: <LazyRoute Component={PublicJobDetailPage} /> },
          { path: ':jobId/apply', element: <LazyRoute Component={PublicApplyPage} /> },
        ],
      },
    ],
  },
])

export function AppRouter() {
  return <RouterProvider router={router} />
}
