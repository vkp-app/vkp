import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
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
};

export type ClusterStatus = {
  __typename?: 'ClusterStatus';
  kubeURL: Scalars['String'];
  kubeVersion: Scalars['String'];
  webURL: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createTenant: Tenant;
};


export type MutationCreateTenantArgs = {
  name: Scalars['String'];
};

export type NamespacedName = {
  __typename?: 'NamespacedName';
  name: Scalars['String'];
  namespace: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  cluster: Cluster;
  clustersInTenant: Array<Cluster>;
  currentUser: User;
  tenants: Array<Tenant>;
};


export type QueryClusterArgs = {
  name: Scalars['ID'];
  tenant: Scalars['ID'];
};


export type QueryClustersInTenantArgs = {
  tenant: Scalars['ID'];
};

export type Tenant = {
  __typename?: 'Tenant';
  name: Scalars['ID'];
  observedClusters: Array<NamespacedName>;
  owner: Scalars['String'];
};

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


export type ClusterQuery = { __typename?: 'Query', cluster: { __typename?: 'Cluster', name: string, tenant: string, status: { __typename?: 'ClusterStatus', kubeVersion: string, kubeURL: string } } };

export type CurrentUserQueryVariables = Exact<{ [key: string]: never; }>;


export type CurrentUserQuery = { __typename?: 'Query', currentUser: { __typename?: 'User', username: string, groups: Array<string> } };


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

      export interface PossibleTypesResultData {
        possibleTypes: {
          [key: string]: string[]
        }
      }
      const result: PossibleTypesResultData = {
  "possibleTypes": {}
};
      export default result;
    