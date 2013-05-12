      struct UserProfile {
        1: i32 uid,
        2: string name,
        3: string blurb
      }
      service UserStorage {
        bool store(1: UserProfile user),
      }