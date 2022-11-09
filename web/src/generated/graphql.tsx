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

export type Cluster = {
  __typename?: 'Cluster';
  name: Scalars['ID'];
  status: ClusterStatus;
  tenant: Scalars['ID'];
  track: Track;
};

export type ClusterStatus = {
  __typename?: 'ClusterStatus';
  kubeURL: Scalars['String'];
  kubeVersion: Scalars['String'];
  webURL: Scalars['String'];
};

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
};


export type MutationApproveTenantArgs = {
  tenant: Scalars['ID'];
};


export type MutationCreateClusterArgs = {
  input: NewCluster;
  tenant: Scalars['ID'];
};


export type MutationCreateTenantArgs = {
  name: Scalars['String'];
};

export type NamespacedName = {
  __typename?: 'NamespacedName';
  name: Scalars['String'];
  namespace: Scalars['String'];
};

export type NewCluster = {
  name: Scalars['String'];
  track: Track;
};

export type Query = {
  __typename?: 'Query';
  cluster: Cluster;
  clusterMetricCPU: Array<MetricValue>;
  clusterMetricMemory: Array<MetricValue>;
  clusterMetricNetReceive: Array<MetricValue>;
  clusterMetricNetTransmit: Array<MetricValue>;
  clusterMetricPods: Array<MetricValue>;
  clustersInTenant: Array<Cluster>;
  currentUser: User;
  tenant: Tenant;
  tenants: Array<Tenant>;
};


export type QueryClusterArgs = {
  name: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClusterMetricCpuArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClusterMetricMemoryArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClusterMetricNetReceiveArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClusterMetricNetTransmitArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClusterMetricPodsArgs = {
  cluster: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClustersInTenantArgs = {
  tenant: Scalars['ID'];
};


export type QueryTenantArgs = {
  name: Scalars['ID'];
};

export type Tenant = {
  __typename?: 'Tenant';
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

export enum Track {
  Beta = 'BETA',
  Rapid = 'RAPID',
  Regular = 'REGULAR',
  Stable = 'STABLE'
}

export type User = {
  __typename?: 'User';
  groups: Array<Scalars['String']>;
  username: Scalars['String'];
};

export type ClustersQueryVariables = Exact<{
  tenant: Scalars['ID'];
}>;


export type ClustersQuery = { __typename?: 'Query', clustersInTenant: Array<{ __typename?: 'Cluster', name: string, tenant: string, status: { __typename?: 'ClusterStatus', kubeVersion: string, kubeURL: string } }> };

export type ClusterQueryVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type ClusterQuery = { __typename?: 'Query', cluster: { __typename?: 'Cluster', name: string, tenant: string, track: Track, status: { __typename?: 'ClusterStatus', kubeVersion: string, kubeURL: string } } };

export type CreateClusterMutationVariables = Exact<{
  tenant: Scalars['ID'];
  input: NewCluster;
}>;


export type CreateClusterMutation = { __typename?: 'Mutation', createCluster: { __typename?: 'Cluster', name: string } };

export type CurrentUserQueryVariables = Exact<{ [key: string]: never; }>;


export type CurrentUserQuery = { __typename?: 'Query', currentUser: { __typename?: 'User', username: string, groups: Array<string> } };

export type MetricsClusterQueryVariables = Exact<{
  tenant: Scalars['ID'];
  cluster: Scalars['ID'];
}>;


export type MetricsClusterQuery = { __typename?: 'Query', clusterMetricMemory: Array<{ __typename?: 'MetricValue', time: number, value: string }>, clusterMetricCPU: Array<{ __typename?: 'MetricValue', time: number, value: string }>, clusterMetricPods: Array<{ __typename?: 'MetricValue', time: number, value: string }>, clusterMetricNetReceive: Array<{ __typename?: 'MetricValue', time: number, value: string }> };

export type TenantsQueryVariables = Exact<{ [key: string]: never; }>;


export type TenantsQuery = { __typename?: 'Query', tenants: Array<{ __typename?: 'Tenant', name: string, owner: string, status: { __typename?: 'TenantStatus', phase: TenantPhase } }> };

export type TenantQueryVariables = Exact<{
  name: Scalars['ID'];
}>;


export type TenantQuery = { __typename?: 'Query', tenant: { __typename?: 'Tenant', name: string, owner: string, status: { __typename?: 'TenantStatus', phase: TenantPhase } } };


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
  cluster(tenant: $tenant, name: $cluster) {
    name
    tenant
    track
    status {
      kubeVersion
      kubeURL
    }
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
export const MetricsClusterDocument = gql`
    query metricsCluster($tenant: ID!, $cluster: ID!) {
  clusterMetricMemory(tenant: $tenant, cluster: $cluster) {
    time
    value
  }
  clusterMetricCPU(tenant: $tenant, cluster: $cluster) {
    time
    value
  }
  clusterMetricPods(tenant: $tenant, cluster: $cluster) {
    time
    value
  }
  clusterMetricNetReceive(tenant: $tenant, cluster: $cluster) {
    time
    value
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
    query tenant($name: ID!) {
  tenant(name: $name) {
    name
    owner
    status {
      phase
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
 *      name: // value for 'name'
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

      export interface PossibleTypesResultData {
        possibleTypes: {
          [key: string]: string[]
        }
      }
      const result: PossibleTypesResultData = {
  "possibleTypes": {}
};
      export default result;
    