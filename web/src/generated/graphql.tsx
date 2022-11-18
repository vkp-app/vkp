import * as Apollo from '@apollo/client';
import {gql} from '@apollo/client';

export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type AccessRef = {
  __typename?: 'AccessRef';
  group: Scalars['String'];
  readOnly: Scalars['Boolean'];
  user: Scalars['String'];
};

export type AccessRefInput = {
  group: Scalars['String'];
  readOnly: Scalars['Boolean'];
  user: Scalars['String'];
};

export type AddonBindingStatus = {
  __typename?: 'AddonBindingStatus';
  name: Scalars['String'];
  phase: AddonPhase;
};

export enum AddonPhase {
  Deleting = 'Deleting',
  Installed = 'Installed',
  Installing = 'Installing'
}

export enum AddonSource {
  Community = 'Community',
  Official = 'Official',
  Platform = 'Platform',
  Unknown = 'Unknown'
}

export type Cluster = {
  __typename?: 'Cluster';
  accessors: Array<AccessRef>;
  name: Scalars['ID'];
  status: ClusterStatus;
  tenant: Scalars['ID'];
  track: ReleaseTrack;
};

export type ClusterAddon = {
  __typename?: 'ClusterAddon';
  description: Scalars['String'];
  displayName: Scalars['String'];
  logo: Scalars['String'];
  maintainer: Scalars['String'];
  name: Scalars['String'];
  source: AddonSource;
  sourceURL: Scalars['String'];
};

export type ClusterStatus = {
  __typename?: 'ClusterStatus';
  kubeURL: Scalars['String'];
  kubeVersion: Scalars['String'];
  platformVersion: Scalars['String'];
  webURL: Scalars['String'];
};

export type Metric = {
  __typename?: 'Metric';
  format: MetricFormat;
  metric: Scalars['String'];
  name: Scalars['String'];
  values: Array<MetricValue>;
};

export enum MetricFormat {
  Bytes = 'Bytes',
  Cpu = 'CPU',
  Plain = 'Plain',
  Time = 'Time'
}

export type MetricValue = {
  __typename?: 'MetricValue';
  time: Scalars['Int'];
  value: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  approveTenant: Scalars['Boolean'];
  createCluster: Cluster;
  createTenant: Tenant;
  deleteCluster: Scalars['Boolean'];
  installAddon: Scalars['Boolean'];
  setClusterAccessors: Scalars['Boolean'];
  setTenantAccessors: Scalars['Boolean'];
  uninstallAddon: Scalars['Boolean'];
};


export type MutationApproveTenantArgs = {
  tenant: Scalars['ID'];
};


export type MutationCreateClusterArgs = {
  input: NewCluster;
  tenant: Scalars['ID'];
};


export type MutationCreateTenantArgs = {
  tenant: Scalars['String'];
};


export type MutationDeleteClusterArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type MutationInstallAddonArgs = {
  addon: Scalars['String'];
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type MutationSetClusterAccessorsArgs = {
  accessors: Array<AccessRefInput>;
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type MutationSetTenantAccessorsArgs = {
  accessors: Array<AccessRefInput>;
  tenant: Scalars['ID'];
};


export type MutationUninstallAddonArgs = {
  addon: Scalars['String'];
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};

export type NamespacedName = {
  __typename?: 'NamespacedName';
  name: Scalars['String'];
  namespace: Scalars['String'];
};

export type NewCluster = {
  ha: Scalars['Boolean'];
  name: Scalars['String'];
  track: ReleaseTrack;
};

export type Query = {
  __typename?: 'Query';
  cluster: Cluster;
  clusterAddons: Array<ClusterAddon>;
  clusterInstalledAddons: Array<AddonBindingStatus>;
  clusterMetrics: Array<Metric>;
  clustersInTenant: Array<Cluster>;
  currentUser: User;
  hasClusterAccess: Scalars['Boolean'];
  hasRole: Scalars['Boolean'];
  hasTenantAccess: Scalars['Boolean'];
  renderKubeconfig: Scalars['String'];
  tenant: Tenant;
  tenants: Array<Tenant>;
};


export type QueryClusterArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClusterAddonsArgs = {
  tenant: Scalars['ID'];
};


export type QueryClusterInstalledAddonsArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClusterMetricsArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClustersInTenantArgs = {
  tenant: Scalars['ID'];
};


export type QueryHasClusterAccessArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
  write: Scalars['Boolean'];
};


export type QueryHasRoleArgs = {
  role: Role;
};


export type QueryHasTenantAccessArgs = {
  tenant: Scalars['ID'];
  write: Scalars['Boolean'];
};


export type QueryRenderKubeconfigArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryTenantArgs = {
  tenant: Scalars['ID'];
};

export enum ReleaseTrack {
  Beta = 'Beta',
  Rapid = 'Rapid',
  Regular = 'Regular',
  Stable = 'Stable'
}

export enum Role {
  Admin = 'ADMIN',
  User = 'USER'
}

export type Tenant = {
  __typename?: 'Tenant';
  accessors: Array<AccessRef>;
  name: Scalars['ID'];
  observedClusters: Array<NamespacedName>;
  owner: Scalars['String'];
  status: TenantStatus;
};

export enum TenantPhase {
  PendingApproval = 'PendingApproval',
  Ready = 'Ready'
}

export type TenantStatus = {
  __typename?: 'TenantStatus';
  phase: TenantPhase;
};

export type User = {
  __typename?: 'User';
  groups: Array<Scalars['String']>;
  username: Scalars['String'];
};

export type TenantAccessQueryVariables = Exact<{
  tenant: Scalars['ID'];
}>;


export type TenantAccessQuery = { __typename?: 'Query', hasTenantAccess: boolean, tenant: { __typename?: 'Tenant', name: string, owner: string, status: { __typename?: 'TenantStatus', phase: TenantPhase }, accessors: Array<{ __typename?: 'AccessRef', user: string, group: string, readOnly: boolean }> } };

export type AllAddonsQueryVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type AllAddonsQuery = { __typename?: 'Query', clusterAddons: Array<{ __typename?: 'ClusterAddon', name: string, displayName: string, maintainer: string, source: AddonSource, sourceURL: string, description: string, logo: string }>, clusterInstalledAddons: Array<{ __typename?: 'AddonBindingStatus', name: string, phase: AddonPhase }> };

export type InstallAddonMutationVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
  addon: Scalars['String'];
}>;


export type InstallAddonMutation = { __typename?: 'Mutation', installAddon: boolean };

export type UninstallAddonMutationVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
  addon: Scalars['String'];
}>;


export type UninstallAddonMutation = { __typename?: 'Mutation', uninstallAddon: boolean };

export type ClusterListQueryVariables = Exact<{
  tenant: Scalars['ID'];
}>;


export type ClusterListQuery = { __typename?: 'Query', hasTenantAccess: boolean, hasRole: boolean, clustersInTenant: Array<{ __typename?: 'Cluster', name: string, tenant: string, status: { __typename?: 'ClusterStatus', kubeVersion: string, kubeURL: string } }>, tenant: { __typename?: 'Tenant', name: string, owner: string, status: { __typename?: 'TenantStatus', phase: TenantPhase } } };

export type ClustersQueryVariables = Exact<{
  tenant: Scalars['ID'];
}>;


export type ClustersQuery = { __typename?: 'Query', clustersInTenant: Array<{ __typename?: 'Cluster', name: string, tenant: string, status: { __typename?: 'ClusterStatus', kubeVersion: string, kubeURL: string } }> };

export type ClusterQueryVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type ClusterQuery = { __typename?: 'Query', cluster: { __typename?: 'Cluster', name: string, tenant: string, track: ReleaseTrack, status: { __typename?: 'ClusterStatus', kubeVersion: string, platformVersion: string, kubeURL: string, webURL: string }, accessors: Array<{ __typename?: 'AccessRef', user: string, group: string, readOnly: boolean }> }, clusterInstalledAddons: Array<{ __typename?: 'AddonBindingStatus', phase: AddonPhase, name: string }> };

export type CreateClusterMutationVariables = Exact<{
  tenant: Scalars['ID'];
  input: NewCluster;
}>;


export type CreateClusterMutation = { __typename?: 'Mutation', createCluster: { __typename?: 'Cluster', name: string } };

export type DeleteClusterMutationVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type DeleteClusterMutation = { __typename?: 'Mutation', deleteCluster: boolean };

export type SetClusterAccessorsMutationVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
  accessors: Array<AccessRefInput> | AccessRefInput;
}>;


export type SetClusterAccessorsMutation = { __typename?: 'Mutation', setClusterAccessors: boolean };

export type KubeConfigQueryVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type KubeConfigQuery = { __typename?: 'Query', renderKubeconfig: string };

export type CurrentUserQueryVariables = Exact<{ [key: string]: never; }>;


export type CurrentUserQuery = { __typename?: 'Query', currentUser: { __typename?: 'User', username: string, groups: Array<string> } };

export type HasAdminQueryVariables = Exact<{ [key: string]: never; }>;


export type HasAdminQuery = { __typename?: 'Query', hasRole: boolean };

export type CanCreateClusterQueryVariables = Exact<{
  tenant: Scalars['ID'];
}>;


export type CanCreateClusterQuery = { __typename?: 'Query', hasTenantAccess: boolean };

export type CanEditClusterQueryVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type CanEditClusterQuery = { __typename?: 'Query', hasClusterAccess: boolean };

export type MetricsClusterQueryVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type MetricsClusterQuery = { __typename?: 'Query', clusterMetrics: Array<{ __typename?: 'Metric', name: string, metric: string, format: MetricFormat, values: Array<{ __typename?: 'MetricValue', value: string }> }> };

export type TenantsQueryVariables = Exact<{ [key: string]: never; }>;


export type TenantsQuery = { __typename?: 'Query', tenants: Array<{ __typename?: 'Tenant', name: string, owner: string, status: { __typename?: 'TenantStatus', phase: TenantPhase } }> };

export type TenantQueryVariables = Exact<{
  tenant: Scalars['ID'];
}>;


export type TenantQuery = { __typename?: 'Query', tenant: { __typename?: 'Tenant', name: string, owner: string, status: { __typename?: 'TenantStatus', phase: TenantPhase }, accessors: Array<{ __typename?: 'AccessRef', user: string, group: string, readOnly: boolean }> } };

export type ApproveTenancyMutationVariables = Exact<{
  tenant: Scalars['ID'];
}>;


export type ApproveTenancyMutation = { __typename?: 'Mutation', approveTenant: boolean };

export type CreateTenantMutationVariables = Exact<{
  tenant: Scalars['String'];
}>;


export type CreateTenantMutation = { __typename?: 'Mutation', createTenant: { __typename?: 'Tenant', name: string } };

export type SetTenantAccessorsMutationVariables = Exact<{
  tenant: Scalars['ID'];
  accessors: Array<AccessRefInput> | AccessRefInput;
}>;


export type SetTenantAccessorsMutation = { __typename?: 'Mutation', setTenantAccessors: boolean };


export const TenantAccessDocument = gql`
    query tenantAccess($tenant: ID!) {
  tenant(tenant: $tenant) {
    name
    owner
    status {
      phase
    }
    accessors {
      user
      group
      readOnly
    }
  }
  hasTenantAccess(tenant: $tenant, write: true)
}
    `;

/**
 * __useTenantAccessQuery__
 *
 * To run a query within a React component, call `useTenantAccessQuery` and pass it any options that fit your needs.
 * When your component renders, `useTenantAccessQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useTenantAccessQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *   },
 * });
 */
export function useTenantAccessQuery(baseOptions: Apollo.QueryHookOptions<TenantAccessQuery, TenantAccessQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<TenantAccessQuery, TenantAccessQueryVariables>(TenantAccessDocument, options);
      }
export function useTenantAccessLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<TenantAccessQuery, TenantAccessQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<TenantAccessQuery, TenantAccessQueryVariables>(TenantAccessDocument, options);
        }
export type TenantAccessQueryHookResult = ReturnType<typeof useTenantAccessQuery>;
export type TenantAccessLazyQueryHookResult = ReturnType<typeof useTenantAccessLazyQuery>;
export type TenantAccessQueryResult = Apollo.QueryResult<TenantAccessQuery, TenantAccessQueryVariables>;
export const AllAddonsDocument = gql`
    query allAddons($tenant: ID!, $cluster: ID!) {
  clusterAddons(tenant: $tenant) {
    name
    displayName
    maintainer
    source
    sourceURL
    description
    logo
  }
  clusterInstalledAddons(tenant: $tenant, cluster: $cluster) {
    name
    phase
  }
}
    `;

/**
 * __useAllAddonsQuery__
 *
 * To run a query within a React component, call `useAllAddonsQuery` and pass it any options that fit your needs.
 * When your component renders, `useAllAddonsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useAllAddonsQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *   },
 * });
 */
export function useAllAddonsQuery(baseOptions: Apollo.QueryHookOptions<AllAddonsQuery, AllAddonsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<AllAddonsQuery, AllAddonsQueryVariables>(AllAddonsDocument, options);
      }
export function useAllAddonsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<AllAddonsQuery, AllAddonsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<AllAddonsQuery, AllAddonsQueryVariables>(AllAddonsDocument, options);
        }
export type AllAddonsQueryHookResult = ReturnType<typeof useAllAddonsQuery>;
export type AllAddonsLazyQueryHookResult = ReturnType<typeof useAllAddonsLazyQuery>;
export type AllAddonsQueryResult = Apollo.QueryResult<AllAddonsQuery, AllAddonsQueryVariables>;
export const InstallAddonDocument = gql`
    mutation installAddon($tenant: ID!, $cluster: ID!, $addon: String!) {
  installAddon(tenant: $tenant, cluster: $cluster, addon: $addon)
}
    `;
export type InstallAddonMutationFn = Apollo.MutationFunction<InstallAddonMutation, InstallAddonMutationVariables>;

/**
 * __useInstallAddonMutation__
 *
 * To run a mutation, you first call `useInstallAddonMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useInstallAddonMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [installAddonMutation, { data, loading, error }] = useInstallAddonMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *      addon: // value for 'addon'
 *   },
 * });
 */
export function useInstallAddonMutation(baseOptions?: Apollo.MutationHookOptions<InstallAddonMutation, InstallAddonMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<InstallAddonMutation, InstallAddonMutationVariables>(InstallAddonDocument, options);
      }
export type InstallAddonMutationHookResult = ReturnType<typeof useInstallAddonMutation>;
export type InstallAddonMutationResult = Apollo.MutationResult<InstallAddonMutation>;
export type InstallAddonMutationOptions = Apollo.BaseMutationOptions<InstallAddonMutation, InstallAddonMutationVariables>;
export const UninstallAddonDocument = gql`
    mutation uninstallAddon($tenant: ID!, $cluster: ID!, $addon: String!) {
  uninstallAddon(tenant: $tenant, cluster: $cluster, addon: $addon)
}
    `;
export type UninstallAddonMutationFn = Apollo.MutationFunction<UninstallAddonMutation, UninstallAddonMutationVariables>;

/**
 * __useUninstallAddonMutation__
 *
 * To run a mutation, you first call `useUninstallAddonMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUninstallAddonMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [uninstallAddonMutation, { data, loading, error }] = useUninstallAddonMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *      addon: // value for 'addon'
 *   },
 * });
 */
export function useUninstallAddonMutation(baseOptions?: Apollo.MutationHookOptions<UninstallAddonMutation, UninstallAddonMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UninstallAddonMutation, UninstallAddonMutationVariables>(UninstallAddonDocument, options);
      }
export type UninstallAddonMutationHookResult = ReturnType<typeof useUninstallAddonMutation>;
export type UninstallAddonMutationResult = Apollo.MutationResult<UninstallAddonMutation>;
export type UninstallAddonMutationOptions = Apollo.BaseMutationOptions<UninstallAddonMutation, UninstallAddonMutationVariables>;
export const ClusterListDocument = gql`
    query clusterList($tenant: ID!) {
  clustersInTenant(tenant: $tenant) {
    name
    tenant
    status {
      kubeVersion
      kubeURL
    }
  }
  tenant(tenant: $tenant) {
    name
    owner
    status {
      phase
    }
  }
  hasTenantAccess(tenant: $tenant, write: true)
  hasRole(role: ADMIN)
}
    `;

/**
 * __useClusterListQuery__
 *
 * To run a query within a React component, call `useClusterListQuery` and pass it any options that fit your needs.
 * When your component renders, `useClusterListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useClusterListQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *   },
 * });
 */
export function useClusterListQuery(baseOptions: Apollo.QueryHookOptions<ClusterListQuery, ClusterListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ClusterListQuery, ClusterListQueryVariables>(ClusterListDocument, options);
      }
export function useClusterListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ClusterListQuery, ClusterListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ClusterListQuery, ClusterListQueryVariables>(ClusterListDocument, options);
        }
export type ClusterListQueryHookResult = ReturnType<typeof useClusterListQuery>;
export type ClusterListLazyQueryHookResult = ReturnType<typeof useClusterListLazyQuery>;
export type ClusterListQueryResult = Apollo.QueryResult<ClusterListQuery, ClusterListQueryVariables>;
export const ClustersDocument = gql`
    query clusters($tenant: ID!) {
  clustersInTenant(tenant: $tenant) {
    name
    tenant
    status {
      kubeVersion
      kubeURL
    }
  }
}
    `;

/**
 * __useClustersQuery__
 *
 * To run a query within a React component, call `useClustersQuery` and pass it any options that fit your needs.
 * When your component renders, `useClustersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useClustersQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *   },
 * });
 */
export function useClustersQuery(baseOptions: Apollo.QueryHookOptions<ClustersQuery, ClustersQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ClustersQuery, ClustersQueryVariables>(ClustersDocument, options);
      }
export function useClustersLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ClustersQuery, ClustersQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ClustersQuery, ClustersQueryVariables>(ClustersDocument, options);
        }
export type ClustersQueryHookResult = ReturnType<typeof useClustersQuery>;
export type ClustersLazyQueryHookResult = ReturnType<typeof useClustersLazyQuery>;
export type ClustersQueryResult = Apollo.QueryResult<ClustersQuery, ClustersQueryVariables>;
export const ClusterDocument = gql`
    query cluster($tenant: ID!, $cluster: ID!) {
  cluster(tenant: $tenant, cluster: $cluster) {
    name
    tenant
    track
    status {
      kubeVersion
      platformVersion
      kubeURL
      webURL
    }
    accessors {
      user
      group
      readOnly
    }
  }
  clusterInstalledAddons(tenant: $tenant, cluster: $cluster) {
    phase
    name
  }
}
    `;

/**
 * __useClusterQuery__
 *
 * To run a query within a React component, call `useClusterQuery` and pass it any options that fit your needs.
 * When your component renders, `useClusterQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useClusterQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *   },
 * });
 */
export function useClusterQuery(baseOptions: Apollo.QueryHookOptions<ClusterQuery, ClusterQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ClusterQuery, ClusterQueryVariables>(ClusterDocument, options);
      }
export function useClusterLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ClusterQuery, ClusterQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ClusterQuery, ClusterQueryVariables>(ClusterDocument, options);
        }
export type ClusterQueryHookResult = ReturnType<typeof useClusterQuery>;
export type ClusterLazyQueryHookResult = ReturnType<typeof useClusterLazyQuery>;
export type ClusterQueryResult = Apollo.QueryResult<ClusterQuery, ClusterQueryVariables>;
export const CreateClusterDocument = gql`
    mutation createCluster($tenant: ID!, $input: NewCluster!) {
  createCluster(tenant: $tenant, input: $input) {
    name
  }
}
    `;
export type CreateClusterMutationFn = Apollo.MutationFunction<CreateClusterMutation, CreateClusterMutationVariables>;

/**
 * __useCreateClusterMutation__
 *
 * To run a mutation, you first call `useCreateClusterMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateClusterMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createClusterMutation, { data, loading, error }] = useCreateClusterMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateClusterMutation(baseOptions?: Apollo.MutationHookOptions<CreateClusterMutation, CreateClusterMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateClusterMutation, CreateClusterMutationVariables>(CreateClusterDocument, options);
      }
export type CreateClusterMutationHookResult = ReturnType<typeof useCreateClusterMutation>;
export type CreateClusterMutationResult = Apollo.MutationResult<CreateClusterMutation>;
export type CreateClusterMutationOptions = Apollo.BaseMutationOptions<CreateClusterMutation, CreateClusterMutationVariables>;
export const DeleteClusterDocument = gql`
    mutation deleteCluster($tenant: ID!, $cluster: ID!) {
  deleteCluster(tenant: $tenant, cluster: $cluster)
}
    `;
export type DeleteClusterMutationFn = Apollo.MutationFunction<DeleteClusterMutation, DeleteClusterMutationVariables>;

/**
 * __useDeleteClusterMutation__
 *
 * To run a mutation, you first call `useDeleteClusterMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteClusterMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteClusterMutation, { data, loading, error }] = useDeleteClusterMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *   },
 * });
 */
export function useDeleteClusterMutation(baseOptions?: Apollo.MutationHookOptions<DeleteClusterMutation, DeleteClusterMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteClusterMutation, DeleteClusterMutationVariables>(DeleteClusterDocument, options);
      }
export type DeleteClusterMutationHookResult = ReturnType<typeof useDeleteClusterMutation>;
export type DeleteClusterMutationResult = Apollo.MutationResult<DeleteClusterMutation>;
export type DeleteClusterMutationOptions = Apollo.BaseMutationOptions<DeleteClusterMutation, DeleteClusterMutationVariables>;
export const SetClusterAccessorsDocument = gql`
    mutation setClusterAccessors($tenant: ID!, $cluster: ID!, $accessors: [AccessRefInput!]!) {
  setClusterAccessors(tenant: $tenant, cluster: $cluster, accessors: $accessors)
}
    `;
export type SetClusterAccessorsMutationFn = Apollo.MutationFunction<SetClusterAccessorsMutation, SetClusterAccessorsMutationVariables>;

/**
 * __useSetClusterAccessorsMutation__
 *
 * To run a mutation, you first call `useSetClusterAccessorsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetClusterAccessorsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setClusterAccessorsMutation, { data, loading, error }] = useSetClusterAccessorsMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *      accessors: // value for 'accessors'
 *   },
 * });
 */
export function useSetClusterAccessorsMutation(baseOptions?: Apollo.MutationHookOptions<SetClusterAccessorsMutation, SetClusterAccessorsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetClusterAccessorsMutation, SetClusterAccessorsMutationVariables>(SetClusterAccessorsDocument, options);
      }
export type SetClusterAccessorsMutationHookResult = ReturnType<typeof useSetClusterAccessorsMutation>;
export type SetClusterAccessorsMutationResult = Apollo.MutationResult<SetClusterAccessorsMutation>;
export type SetClusterAccessorsMutationOptions = Apollo.BaseMutationOptions<SetClusterAccessorsMutation, SetClusterAccessorsMutationVariables>;
export const KubeConfigDocument = gql`
    query kubeConfig($tenant: ID!, $cluster: ID!) {
  renderKubeconfig(tenant: $tenant, cluster: $cluster)
}
    `;

/**
 * __useKubeConfigQuery__
 *
 * To run a query within a React component, call `useKubeConfigQuery` and pass it any options that fit your needs.
 * When your component renders, `useKubeConfigQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useKubeConfigQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *   },
 * });
 */
export function useKubeConfigQuery(baseOptions: Apollo.QueryHookOptions<KubeConfigQuery, KubeConfigQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<KubeConfigQuery, KubeConfigQueryVariables>(KubeConfigDocument, options);
      }
export function useKubeConfigLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<KubeConfigQuery, KubeConfigQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<KubeConfigQuery, KubeConfigQueryVariables>(KubeConfigDocument, options);
        }
export type KubeConfigQueryHookResult = ReturnType<typeof useKubeConfigQuery>;
export type KubeConfigLazyQueryHookResult = ReturnType<typeof useKubeConfigLazyQuery>;
export type KubeConfigQueryResult = Apollo.QueryResult<KubeConfigQuery, KubeConfigQueryVariables>;
export const CurrentUserDocument = gql`
    query currentUser {
  currentUser {
    username
    groups
  }
}
    `;

/**
 * __useCurrentUserQuery__
 *
 * To run a query within a React component, call `useCurrentUserQuery` and pass it any options that fit your needs.
 * When your component renders, `useCurrentUserQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useCurrentUserQuery({
 *   variables: {
 *   },
 * });
 */
export function useCurrentUserQuery(baseOptions?: Apollo.QueryHookOptions<CurrentUserQuery, CurrentUserQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<CurrentUserQuery, CurrentUserQueryVariables>(CurrentUserDocument, options);
      }
export function useCurrentUserLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<CurrentUserQuery, CurrentUserQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<CurrentUserQuery, CurrentUserQueryVariables>(CurrentUserDocument, options);
        }
export type CurrentUserQueryHookResult = ReturnType<typeof useCurrentUserQuery>;
export type CurrentUserLazyQueryHookResult = ReturnType<typeof useCurrentUserLazyQuery>;
export type CurrentUserQueryResult = Apollo.QueryResult<CurrentUserQuery, CurrentUserQueryVariables>;
export const HasAdminDocument = gql`
    query hasAdmin {
  hasRole(role: ADMIN)
}
    `;

/**
 * __useHasAdminQuery__
 *
 * To run a query within a React component, call `useHasAdminQuery` and pass it any options that fit your needs.
 * When your component renders, `useHasAdminQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHasAdminQuery({
 *   variables: {
 *   },
 * });
 */
export function useHasAdminQuery(baseOptions?: Apollo.QueryHookOptions<HasAdminQuery, HasAdminQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<HasAdminQuery, HasAdminQueryVariables>(HasAdminDocument, options);
      }
export function useHasAdminLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<HasAdminQuery, HasAdminQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<HasAdminQuery, HasAdminQueryVariables>(HasAdminDocument, options);
        }
export type HasAdminQueryHookResult = ReturnType<typeof useHasAdminQuery>;
export type HasAdminLazyQueryHookResult = ReturnType<typeof useHasAdminLazyQuery>;
export type HasAdminQueryResult = Apollo.QueryResult<HasAdminQuery, HasAdminQueryVariables>;
export const CanCreateClusterDocument = gql`
    query canCreateCluster($tenant: ID!) {
  hasTenantAccess(tenant: $tenant, write: true)
}
    `;

/**
 * __useCanCreateClusterQuery__
 *
 * To run a query within a React component, call `useCanCreateClusterQuery` and pass it any options that fit your needs.
 * When your component renders, `useCanCreateClusterQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useCanCreateClusterQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *   },
 * });
 */
export function useCanCreateClusterQuery(baseOptions: Apollo.QueryHookOptions<CanCreateClusterQuery, CanCreateClusterQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<CanCreateClusterQuery, CanCreateClusterQueryVariables>(CanCreateClusterDocument, options);
      }
export function useCanCreateClusterLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<CanCreateClusterQuery, CanCreateClusterQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<CanCreateClusterQuery, CanCreateClusterQueryVariables>(CanCreateClusterDocument, options);
        }
export type CanCreateClusterQueryHookResult = ReturnType<typeof useCanCreateClusterQuery>;
export type CanCreateClusterLazyQueryHookResult = ReturnType<typeof useCanCreateClusterLazyQuery>;
export type CanCreateClusterQueryResult = Apollo.QueryResult<CanCreateClusterQuery, CanCreateClusterQueryVariables>;
export const CanEditClusterDocument = gql`
    query canEditCluster($tenant: ID!, $cluster: ID!) {
  hasClusterAccess(tenant: $tenant, cluster: $cluster, write: true)
}
    `;

/**
 * __useCanEditClusterQuery__
 *
 * To run a query within a React component, call `useCanEditClusterQuery` and pass it any options that fit your needs.
 * When your component renders, `useCanEditClusterQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useCanEditClusterQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *   },
 * });
 */
export function useCanEditClusterQuery(baseOptions: Apollo.QueryHookOptions<CanEditClusterQuery, CanEditClusterQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<CanEditClusterQuery, CanEditClusterQueryVariables>(CanEditClusterDocument, options);
      }
export function useCanEditClusterLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<CanEditClusterQuery, CanEditClusterQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<CanEditClusterQuery, CanEditClusterQueryVariables>(CanEditClusterDocument, options);
        }
export type CanEditClusterQueryHookResult = ReturnType<typeof useCanEditClusterQuery>;
export type CanEditClusterLazyQueryHookResult = ReturnType<typeof useCanEditClusterLazyQuery>;
export type CanEditClusterQueryResult = Apollo.QueryResult<CanEditClusterQuery, CanEditClusterQueryVariables>;
export const MetricsClusterDocument = gql`
    query metricsCluster($tenant: ID!, $cluster: ID!) {
  clusterMetrics(tenant: $tenant, cluster: $cluster) {
    name
    metric
    format
    values {
      value
    }
  }
}
    `;

/**
 * __useMetricsClusterQuery__
 *
 * To run a query within a React component, call `useMetricsClusterQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetricsClusterQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetricsClusterQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      cluster: // value for 'cluster'
 *   },
 * });
 */
export function useMetricsClusterQuery(baseOptions: Apollo.QueryHookOptions<MetricsClusterQuery, MetricsClusterQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MetricsClusterQuery, MetricsClusterQueryVariables>(MetricsClusterDocument, options);
      }
export function useMetricsClusterLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MetricsClusterQuery, MetricsClusterQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MetricsClusterQuery, MetricsClusterQueryVariables>(MetricsClusterDocument, options);
        }
export type MetricsClusterQueryHookResult = ReturnType<typeof useMetricsClusterQuery>;
export type MetricsClusterLazyQueryHookResult = ReturnType<typeof useMetricsClusterLazyQuery>;
export type MetricsClusterQueryResult = Apollo.QueryResult<MetricsClusterQuery, MetricsClusterQueryVariables>;
export const TenantsDocument = gql`
    query tenants {
  tenants {
    name
    owner
    status {
      phase
    }
  }
}
    `;

/**
 * __useTenantsQuery__
 *
 * To run a query within a React component, call `useTenantsQuery` and pass it any options that fit your needs.
 * When your component renders, `useTenantsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useTenantsQuery({
 *   variables: {
 *   },
 * });
 */
export function useTenantsQuery(baseOptions?: Apollo.QueryHookOptions<TenantsQuery, TenantsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<TenantsQuery, TenantsQueryVariables>(TenantsDocument, options);
      }
export function useTenantsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<TenantsQuery, TenantsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<TenantsQuery, TenantsQueryVariables>(TenantsDocument, options);
        }
export type TenantsQueryHookResult = ReturnType<typeof useTenantsQuery>;
export type TenantsLazyQueryHookResult = ReturnType<typeof useTenantsLazyQuery>;
export type TenantsQueryResult = Apollo.QueryResult<TenantsQuery, TenantsQueryVariables>;
export const TenantDocument = gql`
    query tenant($tenant: ID!) {
  tenant(tenant: $tenant) {
    name
    owner
    status {
      phase
    }
    accessors {
      user
      group
      readOnly
    }
  }
}
    `;

/**
 * __useTenantQuery__
 *
 * To run a query within a React component, call `useTenantQuery` and pass it any options that fit your needs.
 * When your component renders, `useTenantQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useTenantQuery({
 *   variables: {
 *      tenant: // value for 'tenant'
 *   },
 * });
 */
export function useTenantQuery(baseOptions: Apollo.QueryHookOptions<TenantQuery, TenantQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<TenantQuery, TenantQueryVariables>(TenantDocument, options);
      }
export function useTenantLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<TenantQuery, TenantQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<TenantQuery, TenantQueryVariables>(TenantDocument, options);
        }
export type TenantQueryHookResult = ReturnType<typeof useTenantQuery>;
export type TenantLazyQueryHookResult = ReturnType<typeof useTenantLazyQuery>;
export type TenantQueryResult = Apollo.QueryResult<TenantQuery, TenantQueryVariables>;
export const ApproveTenancyDocument = gql`
    mutation approveTenancy($tenant: ID!) {
  approveTenant(tenant: $tenant)
}
    `;
export type ApproveTenancyMutationFn = Apollo.MutationFunction<ApproveTenancyMutation, ApproveTenancyMutationVariables>;

/**
 * __useApproveTenancyMutation__
 *
 * To run a mutation, you first call `useApproveTenancyMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useApproveTenancyMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [approveTenancyMutation, { data, loading, error }] = useApproveTenancyMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *   },
 * });
 */
export function useApproveTenancyMutation(baseOptions?: Apollo.MutationHookOptions<ApproveTenancyMutation, ApproveTenancyMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ApproveTenancyMutation, ApproveTenancyMutationVariables>(ApproveTenancyDocument, options);
      }
export type ApproveTenancyMutationHookResult = ReturnType<typeof useApproveTenancyMutation>;
export type ApproveTenancyMutationResult = Apollo.MutationResult<ApproveTenancyMutation>;
export type ApproveTenancyMutationOptions = Apollo.BaseMutationOptions<ApproveTenancyMutation, ApproveTenancyMutationVariables>;
export const CreateTenantDocument = gql`
    mutation createTenant($tenant: String!) {
  createTenant(tenant: $tenant) {
    name
  }
}
    `;
export type CreateTenantMutationFn = Apollo.MutationFunction<CreateTenantMutation, CreateTenantMutationVariables>;

/**
 * __useCreateTenantMutation__
 *
 * To run a mutation, you first call `useCreateTenantMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateTenantMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createTenantMutation, { data, loading, error }] = useCreateTenantMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *   },
 * });
 */
export function useCreateTenantMutation(baseOptions?: Apollo.MutationHookOptions<CreateTenantMutation, CreateTenantMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateTenantMutation, CreateTenantMutationVariables>(CreateTenantDocument, options);
      }
export type CreateTenantMutationHookResult = ReturnType<typeof useCreateTenantMutation>;
export type CreateTenantMutationResult = Apollo.MutationResult<CreateTenantMutation>;
export type CreateTenantMutationOptions = Apollo.BaseMutationOptions<CreateTenantMutation, CreateTenantMutationVariables>;
export const SetTenantAccessorsDocument = gql`
    mutation setTenantAccessors($tenant: ID!, $accessors: [AccessRefInput!]!) {
  setTenantAccessors(tenant: $tenant, accessors: $accessors)
}
    `;
export type SetTenantAccessorsMutationFn = Apollo.MutationFunction<SetTenantAccessorsMutation, SetTenantAccessorsMutationVariables>;

/**
 * __useSetTenantAccessorsMutation__
 *
 * To run a mutation, you first call `useSetTenantAccessorsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetTenantAccessorsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setTenantAccessorsMutation, { data, loading, error }] = useSetTenantAccessorsMutation({
 *   variables: {
 *      tenant: // value for 'tenant'
 *      accessors: // value for 'accessors'
 *   },
 * });
 */
export function useSetTenantAccessorsMutation(baseOptions?: Apollo.MutationHookOptions<SetTenantAccessorsMutation, SetTenantAccessorsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetTenantAccessorsMutation, SetTenantAccessorsMutationVariables>(SetTenantAccessorsDocument, options);
      }
export type SetTenantAccessorsMutationHookResult = ReturnType<typeof useSetTenantAccessorsMutation>;
export type SetTenantAccessorsMutationResult = Apollo.MutationResult<SetTenantAccessorsMutation>;
export type SetTenantAccessorsMutationOptions = Apollo.BaseMutationOptions<SetTenantAccessorsMutation, SetTenantAccessorsMutationVariables>;

      export interface PossibleTypesResultData {
        possibleTypes: {
          [key: string]: string[]
        }
      }
      const result: PossibleTypesResultData = {
  "possibleTypes": {}
};
      export default result;
    